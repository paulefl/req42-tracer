package aspice

import (
	"testing"

	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
)

func buildCheckerGraph() *graph.Analyzer {
	g := &model.TraceabilityGraph{
		Requirements: map[string]*model.Requirement{
			"REQ-001": {ID: "REQ-001"},
			"REQ-002": {ID: "REQ-002"},
		},
		ArchElements: map[string]*model.ArchElement{
			"comp.api":    {ID: "comp.api", Parent: "comp.system", Req: []string{"REQ-001"}, Impl: "internal/api.go"},
			"comp.system": {ID: "comp.system", Req: []string{"REQ-002"}},
		},
		TestSpecs: map[string]*model.TestSpec{
			"TS-001": {ID: "TS-001", Req: []string{"REQ-001"}},
		},
		TestCodes:   make(map[string]*model.TestCode),
		TestResults: make(map[string]*model.TestResult),
		Links: []*model.TraceLink{
			{FromID: "REQ-001", ToID: "comp.api", LinkType: "satisfied-by", Status: "active"},
			{FromID: "REQ-001", ToID: "TS-001", LinkType: "verified-by", Status: "active"},
		},
	}
	return graph.NewAnalyzer(g)
}

func fullConfig() *model.Config {
	return aspiceConfig("SWE.1", "SWE.2", "SWE.3", "SWE.5")
}

// aspiceConfig builds a Config with the given ASPICE process IDs — avoids
// repeating the verbose anonymous-struct literal in every test function.
func aspiceConfig(processes ...string) *model.Config {
	return &model.Config{
		ASPICE: struct {
			AutoDerive   bool                         `yaml:"auto-derive"`
			Processes    []string                     `yaml:"processes"`
			ProcessRules map[string]map[string]string `yaml:"process-rules"`
		}{
			Processes: processes,
		},
	}
}

// [test-spec,id=TS-ASPICE-001,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestNewChecker verifies that NewChecker initializes correctly.
func TestNewChecker(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), fullConfig())
	if checker == nil {
		t.Fatal("NewChecker returned nil")
	}
}

// [test-spec,id=TS-ASPICE-002,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestCheckCompliance verifies that CheckCompliance returns results for all configured processes.
func TestCheckCompliance(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), fullConfig())
	report := checker.CheckCompliance()
	if report == nil {
		t.Fatal("CheckCompliance returned nil")
	}
	if len(report.Processes) == 0 {
		t.Error("expected processes in report")
	}
	// All 4 configured processes should be present
	for _, pid := range []string{"SWE.1", "SWE.2", "SWE.3", "SWE.5"} {
		if _, ok := report.Processes[pid]; !ok {
			t.Errorf("process %s missing from report", pid)
		}
	}
}

// [test-spec,id=TS-ASPICE-003,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestCheckCompliance_DefaultProcesses verifies default process selection when config is empty.
func TestCheckCompliance_DefaultProcesses(t *testing.T) {
	config := &model.Config{}
	checker := NewChecker(buildCheckerGraph(), config)
	report := checker.CheckCompliance()
	if len(report.Processes) == 0 {
		t.Error("expected default processes in report")
	}
}

// [test-spec,id=TS-ASPICE-004,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestCheckCompliance_Overall verifies overall score is between 0 and 100.
func TestCheckCompliance_Overall(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), fullConfig())
	report := checker.CheckCompliance()
	if report.Overall < 0 || report.Overall > 100 {
		t.Errorf("Overall coverage out of range: %.2f", report.Overall)
	}
}

// [test-spec,id=TS-ASPICE-005,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestGetProcessCoverage_SWE1 verifies coverage for SWE.1 process.
func TestGetProcessCoverage_SWE1(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), fullConfig())
	cov, err := checker.GetProcessCoverage("SWE.1")
	if err != nil {
		t.Fatalf("GetProcessCoverage error: %v", err)
	}
	if cov < 0 || cov > 100 {
		t.Errorf("SWE.1 coverage out of range: %.1f", cov)
	}
}

// [test-spec,id=TS-ASPICE-006,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestGetProcessCoverage_UnknownProcess verifies error for unknown process ID.
func TestGetProcessCoverage_UnknownProcess(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), fullConfig())
	_, err := checker.GetProcessCoverage("UNKNOWN.99")
	if err == nil {
		t.Error("expected error for unknown process")
	}
}

// [test-spec,id=TS-ASPICE-007,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestCheckCompliance_UnknownProcessSkipped verifies that unknown processes are skipped.
func TestCheckCompliance_UnknownProcessSkipped(t *testing.T) {
	config := aspiceConfig("UNKNOWN.99")
	checker := NewChecker(buildCheckerGraph(), config)
	report := checker.CheckCompliance()
	if report == nil {
		t.Fatal("expected non-nil report even with unknown processes")
	}
	if len(report.Processes) != 0 {
		t.Errorf("expected 0 processes for all-unknown config, got %d", len(report.Processes))
	}
}

// [test-spec,id=TS-ASPICE-008,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestSWE1_BP2_PartialCoverage verifies SWE.1.BP2 is partial when not all reqs have arch links.
func TestSWE1_BP2_PartialCoverage(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), aspiceConfig("SWE.1"))
	report := checker.CheckCompliance()
	results := report.Processes["SWE.1"]
	var bp2 *model.ASPICECheckResult
	for _, r := range results {
		if r.BP.ID == "SWE.1.BP2" {
			bp2 = r
		}
	}
	if bp2 == nil {
		t.Fatal("SWE.1.BP2 result not found")
	}
	// REQ-002 is not in a satisfied-by link → partial
	if bp2.Coverage == 100 {
		t.Error("expected partial coverage for SWE.1.BP2 since REQ-002 has no arch link")
	}
}

// [test-spec,id=TS-ASPICE-009,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestSWE1_BP6_TestCoverage verifies SWE.1.BP6 testability check.
func TestSWE1_BP6_TestCoverage(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), aspiceConfig("SWE.1"))
	report := checker.CheckCompliance()
	results := report.Processes["SWE.1"]
	var bp6 *model.ASPICECheckResult
	for _, r := range results {
		if r.BP.ID == "SWE.1.BP6" {
			bp6 = r
		}
	}
	if bp6 == nil {
		t.Fatal("SWE.1.BP6 not found")
	}
	if bp6.Coverage < 0 || bp6.Coverage > 100 {
		t.Errorf("SWE.1.BP6 coverage out of range: %.1f", bp6.Coverage)
	}
}

// [test-spec,id=TS-ASPICE-010,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestSWE2_BP4_Coverage verifies SWE.2.BP4 arch traceability check.
func TestSWE2_BP4_Coverage(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), aspiceConfig("SWE.2"))
	report := checker.CheckCompliance()
	results := report.Processes["SWE.2"]
	if len(results) == 0 {
		t.Fatal("no SWE.2 results")
	}
	var bp4 *model.ASPICECheckResult
	for _, r := range results {
		if r.BP.ID == "SWE.2.BP4" {
			bp4 = r
		}
	}
	if bp4 == nil {
		t.Fatal("SWE.2.BP4 not found")
	}
	// Both arch elements have req → 100%
	if bp4.Coverage != 100 {
		t.Errorf("SWE.2.BP4 coverage = %.1f, want 100", bp4.Coverage)
	}
}

// [test-spec,id=TS-ASPICE-011,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestSWE3_BP3_ImplCoverage verifies SWE.3.BP3 implementation traceability check.
func TestSWE3_BP3_ImplCoverage(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), aspiceConfig("SWE.3"))
	report := checker.CheckCompliance()
	var bp3 *model.ASPICECheckResult
	for _, r := range report.Processes["SWE.3"] {
		if r.BP.ID == "SWE.3.BP3" {
			bp3 = r
		}
	}
	if bp3 == nil {
		t.Fatal("SWE.3.BP3 not found")
	}
	// comp.api has impl, comp.system has no parent → only comp.api counted
	if bp3.Coverage != 100 {
		t.Errorf("SWE.3.BP3 coverage = %.1f, want 100 (comp.api has impl)", bp3.Coverage)
	}
}

// [test-spec,id=TS-ASPICE-012,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestSWE5_BP3_TestTraceability verifies SWE.5.BP3 test-to-req traceability check.
func TestSWE5_BP3_TestTraceability(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), aspiceConfig("SWE.5"))
	report := checker.CheckCompliance()
	var bp3 *model.ASPICECheckResult
	for _, r := range report.Processes["SWE.5"] {
		if r.BP.ID == "SWE.5.BP3" {
			bp3 = r
		}
	}
	if bp3 == nil {
		t.Fatal("SWE.5.BP3 not found")
	}
	// TS-001 has req=REQ-001 → 100%
	if bp3.Coverage != 100 {
		t.Errorf("SWE.5.BP3 coverage = %.1f, want 100", bp3.Coverage)
	}
}

// [test-spec,id=TS-ASPICE-013,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestCheckCompliance_EmptyGraph verifies compliance with empty graph returns zero coverage.
func TestCheckCompliance_EmptyGraph(t *testing.T) {
	g := &model.TraceabilityGraph{
		Requirements: make(map[string]*model.Requirement),
		ArchElements: make(map[string]*model.ArchElement),
		TestSpecs:    make(map[string]*model.TestSpec),
		TestCodes:    make(map[string]*model.TestCode),
		TestResults:  make(map[string]*model.TestResult),
		Links:        []*model.TraceLink{},
	}
	analyzer := graph.NewAnalyzer(g)
	checker := NewChecker(analyzer, fullConfig())
	report := checker.CheckCompliance()
	if report.Overall != 0 {
		t.Errorf("expected 0 overall for empty graph, got %.2f", report.Overall)
	}
}

// [test-spec,id=TS-ASPICE-014,req="REQ-ASPICE-001",aspice="SWE.5.BP3"]
// TestGetProcessCoverage_AllProcesses verifies coverage is calculable for all registered processes.
func TestGetProcessCoverage_AllProcesses(t *testing.T) {
	checker := NewChecker(buildCheckerGraph(), fullConfig())
	for _, pid := range []string{"SWE.1", "SWE.2", "SWE.3", "SWE.5"} {
		cov, err := checker.GetProcessCoverage(pid)
		if err != nil {
			t.Errorf("GetProcessCoverage(%q) error: %v", pid, err)
		}
		if cov < 0 || cov > 100 {
			t.Errorf("GetProcessCoverage(%q) = %.1f out of range", pid, cov)
		}
	}
}
