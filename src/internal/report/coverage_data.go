package report

import (
	"sort"
	"strings"

	"github.com/paulefl/req42-tracer/src/internal/model"
	"github.com/paulefl/req42-tracer/src/internal/testresult"
)

// CoverageRow is one row in the Coverage-Tab table.
type CoverageRow struct {
	ArchID     string  `json:"arch_id"`     // matched arch element ID (empty if unmatched)
	ArchTitle  string  `json:"arch_title"`
	Package    string  `json:"package"`
	Statements int     `json:"statements"`
	Covered    int     `json:"covered"`
	Pct        float64 `json:"pct"`
	Level      string  `json:"level"` // good / warning / danger
}

// CoverageData is injected into the Coverage-Tab.
type CoverageData struct {
	Rows        []CoverageRow `json:"rows"`
	TotalStmts  int           `json:"total_stmts"`
	TotalCov    int           `json:"total_covered"`
	OverallPct  float64       `json:"overall_pct"`
	OverallLevel string       `json:"overall_level"`
}

// BuildCoverageData merges PackageCoverage entries with the traceability graph.
// It matches each package to an arch element via the impl= field.
func BuildCoverageData(pkgs []testresult.PackageCoverage, g *model.TraceabilityGraph) *CoverageData {
	// Build impl→archID lookup (last segment of impl path → arch)
	implIndex := make(map[string]*model.ArchElement)
	for _, arch := range g.ArchElements {
		if arch.Impl == "" {
			continue
		}
		short := arch.Impl
		if i := strings.LastIndex(short, "/"); i >= 0 {
			short = short[i+1:]
		}
		implIndex[strings.ToLower(short)] = arch
		// Also index the full impl path
		implIndex[strings.ToLower(arch.Impl)] = arch
	}

	data := &CoverageData{Rows: make([]CoverageRow, 0, len(pkgs))}

	for _, pkg := range pkgs {
		row := CoverageRow{
			Package:    pkg.Package,
			Statements: pkg.Statements,
			Covered:    pkg.Covered,
			Pct:        pkg.Pct,
			Level:      coverageLevel(pkg.Pct),
		}
		// Try to match arch element
		key := strings.ToLower(pkg.Package)
		if arch, ok := implIndex[key]; ok {
			row.ArchID = arch.ID
			row.ArchTitle = arch.Title
		}
		data.Rows = append(data.Rows, row)
		data.TotalStmts += pkg.Statements
		data.TotalCov += pkg.Covered
	}

	// Sort: unmatched first (gaps), then by pct ascending (critical first)
	sort.Slice(data.Rows, func(i, j int) bool {
		if data.Rows[i].ArchID == "" && data.Rows[j].ArchID != "" {
			return true
		}
		if data.Rows[i].ArchID != "" && data.Rows[j].ArchID == "" {
			return false
		}
		return data.Rows[i].Pct < data.Rows[j].Pct
	})

	if data.TotalStmts > 0 {
		data.OverallPct = float64(data.TotalCov) / float64(data.TotalStmts) * 100
	}
	data.OverallLevel = coverageLevel(data.OverallPct)
	return data
}

func coverageLevel(pct float64) string {
	if pct >= 80 {
		return "good"
	}
	if pct >= 70 {
		return "warning"
	}
	return "danger"
}
