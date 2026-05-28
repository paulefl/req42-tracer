package report

// [test-spec,id=TS-RPT-001,req="REQ-RPT-010",aspice="SWE.5.BP3"]
// Test: BuildASPICEDashboardData mit minimalen Graph-Daten
// Verifies that ASPICE dashboard data is built correctly from a traceability graph.
// [end]

// [test-spec,id=TS-RPT-002,req="REQ-RPT-010",aspice="SWE.5.BP3"]
// Test: BuildASPICEDashboardData mit leeren Prozesslisten
// Verifies that empty process configuration falls back to default processes.
// [end]

// [test-spec,id=TS-RPT-003,req="REQ-RPT-011",aspice="SWE.5.BP3"]
// Test: getCoverageLevel Schwellenwerte
// Verifies that coverage thresholds (good/warning/danger) are applied correctly.
// [end]

import (
	"testing"

	"github.com/paulefl/req42-tracer/internal/graph"
	"github.com/paulefl/req42-tracer/internal/model"
)

func makeTestConfig(processes []string) *model.Config {
	if len(processes) == 0 {
		processes = []string{"SWE.1", "SWE.2", "SWE.3", "SWE.5"}
	}
	return &model.Config{
		ASPICE: struct {
			AutoDerive   bool              `yaml:"auto-derive"`
			Processes    []string          `yaml:"processes"`
			ProcessRules map[string]map[string]string `yaml:"process-rules"`
		}{
			Processes:    processes,
			ProcessRules: make(map[string]map[string]string),
		},
	}
}

func makeTestGraph() *model.TraceabilityGraph {
	return &model.TraceabilityGraph{
		Requirements: map[string]*model.Requirement{
			"REQ-001": {ID: "REQ-001", Title: "Req 1", ASPICE: "SWE.1"},
			"REQ-002": {ID: "REQ-002", Title: "Req 2", ASPICE: "SWE.1"},
		},
		ArchElements: map[string]*model.ArchElement{
			"comp.api": {ID: "comp.api", Title: "API", ASPICE: "SWE.2", Req: []string{"REQ-001"}, Parent: "system"},
		},
		TestSpecs: map[string]*model.TestSpec{
			"spec.api": {ID: "spec.api", Title: "API Spec", Req: []string{"REQ-001"}},
		},
		TestResults: make(map[string]*model.TestResult),
		TestCodes:   make(map[string]*model.TestCode),
		Links: []*model.TraceLink{
			{FromID: "REQ-001", FromType: "requirement", ToID: "comp.api", ToType: "arch", LinkType: "satisfied-by"},
			{FromID: "REQ-001", FromType: "requirement", ToID: "spec.api", ToType: "test-spec", LinkType: "verified-by"},
			{FromID: "spec.api", FromType: "test-spec", ToID: "REQ-001", ToType: "requirement", LinkType: "verifies"},
		},
	}
}

// TS-RPT-001: BuildASPICEDashboardData mit minimalen Graph-Daten
func TestBuildASPICEDashboardData_Basic(t *testing.T) {
	g := makeTestGraph()
	analyzer := graph.NewAnalyzer(g)
	config := makeTestConfig([]string{"SWE.1", "SWE.2"})

	data := BuildASPICEDashboardData(analyzer, config)

	if data == nil {
		t.Fatal("Expected non-nil dashboard data")
	}

	if len(data.Processes) != 2 {
		t.Errorf("Expected 2 processes, got %d", len(data.Processes))
	}

	// Overall should be between 0 and 100
	if data.Overall < 0 || data.Overall > 100 {
		t.Errorf("Overall coverage %f out of range [0,100]", data.Overall)
	}

	// Check process IDs are present
	ids := make(map[string]bool)
	for _, p := range data.Processes {
		ids[p.ID] = true
		if p.Name == "" {
			t.Errorf("Process %s has empty name", p.ID)
		}
		if p.Coverage < 0 || p.Coverage > 100 {
			t.Errorf("Process %s coverage %f out of range", p.ID, p.Coverage)
		}
		if p.Status != "good" && p.Status != "warning" && p.Status != "danger" {
			t.Errorf("Process %s has invalid status %q", p.ID, p.Status)
		}
	}

	if !ids["SWE.1"] {
		t.Error("Expected SWE.1 in processes")
	}
	if !ids["SWE.2"] {
		t.Error("Expected SWE.2 in processes")
	}
}

// TS-RPT-001: BPs sind nach ID sortiert
func TestBuildASPICEDashboardData_BPsSorted(t *testing.T) {
	g := makeTestGraph()
	analyzer := graph.NewAnalyzer(g)
	config := makeTestConfig([]string{"SWE.1"})

	data := BuildASPICEDashboardData(analyzer, config)

	if len(data.Processes) == 0 {
		t.Fatal("Expected at least one process")
	}

	bps := data.Processes[0].BPs
	for i := 1; i < len(bps); i++ {
		if bps[i].ID < bps[i-1].ID {
			t.Errorf("BPs not sorted: %s before %s", bps[i-1].ID, bps[i].ID)
		}
	}
}

// TS-RPT-001: Gaps sind nie nil (immer leeres Slice)
func TestBuildASPICEDashboardData_GapsNeverNil(t *testing.T) {
	g := makeTestGraph()
	analyzer := graph.NewAnalyzer(g)
	config := makeTestConfig([]string{"SWE.1", "SWE.2", "SWE.3", "SWE.5"})

	data := BuildASPICEDashboardData(analyzer, config)

	for _, proc := range data.Processes {
		for _, bp := range proc.BPs {
			if bp.Gaps == nil {
				t.Errorf("BP %s has nil Gaps slice (expected empty slice)", bp.ID)
			}
		}
	}
}

// TS-RPT-002: Leere Prozessliste fällt auf Defaults zurück
func TestBuildASPICEDashboardData_DefaultProcesses(t *testing.T) {
	g := makeTestGraph()
	analyzer := graph.NewAnalyzer(g)

	// Config with no ASPICE processes set — BuildASPICEDashboardData uses default list
	config := makeTestConfig([]string{})

	data := BuildASPICEDashboardData(analyzer, config)

	if len(data.Processes) == 0 {
		t.Error("Expected default processes when config list is empty")
	}
}

// TS-RPT-002: Unbekannter Prozess wird übersprungen (keine Panic)
func TestBuildASPICEDashboardData_UnknownProcessSkipped(t *testing.T) {
	g := makeTestGraph()
	analyzer := graph.NewAnalyzer(g)
	config := makeTestConfig([]string{"SWE.1", "SWE.99"})

	data := BuildASPICEDashboardData(analyzer, config)

	// SWE.99 is unknown and must be silently skipped
	if len(data.Processes) != 1 {
		t.Errorf("Expected 1 process (SWE.99 skipped), got %d", len(data.Processes))
	}
	if data.Processes[0].ID != "SWE.1" {
		t.Errorf("Expected SWE.1, got %s", data.Processes[0].ID)
	}
}

// TS-RPT-002: Leerer Graph ergibt 0% Coverage, keine Panic
func TestBuildASPICEDashboardData_EmptyGraph(t *testing.T) {
	g := &model.TraceabilityGraph{
		Requirements: make(map[string]*model.Requirement),
		ArchElements: make(map[string]*model.ArchElement),
		TestSpecs:    make(map[string]*model.TestSpec),
		TestResults:  make(map[string]*model.TestResult),
		TestCodes:    make(map[string]*model.TestCode),
		Links:        []*model.TraceLink{},
	}
	analyzer := graph.NewAnalyzer(g)
	config := makeTestConfig([]string{"SWE.1"})

	data := BuildASPICEDashboardData(analyzer, config)

	if data == nil {
		t.Fatal("Expected non-nil result for empty graph")
	}
	if data.Overall != 0 {
		t.Errorf("Expected 0%% overall for empty graph, got %f", data.Overall)
	}
}

// TS-RPT-003: getCoverageLevel Schwellenwerte
func TestGetCoverageLevel(t *testing.T) {
	cases := []struct {
		pct      float64
		expected string
	}{
		{100.0, "good"},
		{80.0, "good"},
		{79.9, "warning"},
		{50.0, "warning"},
		{49.9, "danger"},
		{0.0, "danger"},
	}

	for _, tc := range cases {
		got := getCoverageLevel(tc.pct)
		if got != tc.expected {
			t.Errorf("getCoverageLevel(%.1f) = %q, want %q", tc.pct, got, tc.expected)
		}
	}
}
