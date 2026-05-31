package report

import (
	"testing"

	"github.com/paulefl/req42-tracer/src/internal/model"
	"github.com/paulefl/req42-tracer/src/internal/testresult"
)

// [test-spec,id=TS-RPT-055,req="REQ-REPORT-001",aspice="SWE.5.BP3"]
// TestBuildCoverageData_ArchMapping verifies packages are matched to arch elements via impl=.
func TestBuildCoverageData_ArchMapping(t *testing.T) {
	g := &model.TraceabilityGraph{
		Requirements: make(map[string]*model.Requirement),
		ArchElements: map[string]*model.ArchElement{
			"comp.parser": {ID: "comp.parser", Title: "Parser", Impl: "internal/parser"},
		},
		DesignElements: make(map[string]*model.DesignElement),
		TestSpecs:      make(map[string]*model.TestSpec),
		TestCodes:      make(map[string]*model.TestCode),
		TestResults:    make(map[string]*model.TestResult),
		Links:          []*model.TraceLink{},
	}

	pkgs := []testresult.PackageCoverage{
		{Package: "parser", Statements: 100, Covered: 90, Pct: 90.0},
		{Package: "unknown", Statements: 50, Covered: 30, Pct: 60.0},
	}

	data := BuildCoverageData(pkgs, g)

	if len(data.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(data.Rows))
	}

	// Find parser row
	var parserRow *CoverageRow
	for i := range data.Rows {
		if data.Rows[i].Package == "parser" {
			parserRow = &data.Rows[i]
		}
	}
	if parserRow == nil {
		t.Fatal("parser row not found")
	}
	if parserRow.ArchID != "comp.parser" {
		t.Errorf("ArchID = %q, want comp.parser", parserRow.ArchID)
	}
	if parserRow.Level != "good" {
		t.Errorf("Level = %q, want good (90%%)", parserRow.Level)
	}
}

// [test-spec,id=TS-RPT-056,req="REQ-REPORT-001",aspice="SWE.5.BP3"]
// TestBuildCoverageData_Totals verifies overall percentage and totals are calculated.
func TestBuildCoverageData_Totals(t *testing.T) {
	g := &model.TraceabilityGraph{
		Requirements:   make(map[string]*model.Requirement),
		ArchElements:   make(map[string]*model.ArchElement),
		DesignElements: make(map[string]*model.DesignElement),
		TestSpecs:      make(map[string]*model.TestSpec),
		TestCodes:      make(map[string]*model.TestCode),
		TestResults:    make(map[string]*model.TestResult),
		Links:          []*model.TraceLink{},
	}

	pkgs := []testresult.PackageCoverage{
		{Package: "a", Statements: 100, Covered: 80, Pct: 80.0},
		{Package: "b", Statements: 100, Covered: 60, Pct: 60.0},
	}

	data := BuildCoverageData(pkgs, g)

	if data.TotalStmts != 200 {
		t.Errorf("TotalStmts = %d, want 200", data.TotalStmts)
	}
	if data.TotalCov != 140 {
		t.Errorf("TotalCov = %d, want 140", data.TotalCov)
	}
	if data.OverallPct < 69 || data.OverallPct > 71 {
		t.Errorf("OverallPct = %.1f, want ~70.0", data.OverallPct)
	}
	if data.OverallLevel != "warning" {
		t.Errorf("OverallLevel = %q, want warning", data.OverallLevel)
	}
}

// [test-spec,id=TS-RPT-057,req="REQ-REPORT-001",aspice="SWE.5.BP3"]
// TestCoverageLevel verifies the level classification thresholds.
func TestCoverageLevel(t *testing.T) {
	cases := []struct{ pct float64; want string }{
		{90.0, "good"},
		{80.0, "good"},
		{79.9, "warning"},
		{70.0, "warning"},
		{69.9, "danger"},
		{0.0,  "danger"},
	}
	for _, c := range cases {
		got := coverageLevel(c.pct)
		if got != c.want {
			t.Errorf("coverageLevel(%.1f) = %q, want %q", c.pct, got, c.want)
		}
	}
}
