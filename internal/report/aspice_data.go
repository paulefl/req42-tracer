package report

import (
	"sort"

	"github.com/paulefl/req42-tracer/internal/aspice"
	"github.com/paulefl/req42-tracer/internal/graph"
	"github.com/paulefl/req42-tracer/internal/model"
)

// ASPICEDashboardData holds all data for the ASPICE Dashboard tab.
type ASPICEDashboardData struct {
	Overall   float64             `json:"overall"`
	Processes []ASPICEProcessData `json:"processes"`
}

// ASPICEProcessData holds coverage data for one ASPICE process (e.g. SWE.1).
type ASPICEProcessData struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Coverage    float64        `json:"coverage"`
	Status      string         `json:"status"`
	BPs         []ASPICEBPData `json:"bps"`
}

// ASPICEBPData holds coverage data for one best practice within a process.
type ASPICEBPData struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Coverage    float64  `json:"coverage"`
	Status      string   `json:"status"`
	Gaps        []string `json:"gaps"`
}

// BuildASPICEDashboardData runs the ASPICE compliance check and returns dashboard-ready data.
func BuildASPICEDashboardData(analyzer *graph.Analyzer, config *model.Config) *ASPICEDashboardData {
	checker := aspice.NewChecker(analyzer, config)
	registry := aspice.NewProcessRegistry()
	report := checker.CheckCompliance()

	processIDs := config.ASPICE.Processes
	if len(processIDs) == 0 {
		processIDs = []string{"SWE.1", "SWE.2", "SWE.3", "SWE.5"}
	}

	data := &ASPICEDashboardData{
		Overall:   report.Overall,
		Processes: make([]ASPICEProcessData, 0, len(processIDs)),
	}

	for _, pid := range processIDs {
		proc := registry.GetProcess(pid)
		if proc == nil {
			continue
		}
		results := report.Processes[pid]

		bps := make([]ASPICEBPData, 0, len(results))
		totalCov := 0.0
		for _, r := range results {
			gaps := r.Gaps
			if gaps == nil {
				gaps = []string{}
			}
			bps = append(bps, ASPICEBPData{
				ID:          r.BP.ID,
				Title:       r.BP.Title,
				Description: r.BP.Description,
				Coverage:    r.Coverage,
				Status:      r.Status,
				Gaps:        gaps,
			})
			totalCov += r.Coverage
		}

		avgCov := 0.0
		if len(results) > 0 {
			avgCov = totalCov / float64(len(results))
		}

		sort.Slice(bps, func(i, j int) bool { return bps[i].ID < bps[j].ID })

		data.Processes = append(data.Processes, ASPICEProcessData{
			ID:          proc.ID,
			Name:        proc.Name,
			Description: proc.Description,
			Coverage:    avgCov,
			Status:      getCoverageLevel(avgCov),
			BPs:         bps,
		})
	}

	return data
}
