package validation

import (
	"strings"
	"testing"

	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
)

func buildValidationGraph() *graph.Analyzer {
	g := &model.TraceabilityGraph{
		Requirements: map[string]*model.Requirement{
			"REQ-001": {ID: "REQ-001", ASPICE: "SWE.1", Priority: "high", ReviewedBy: "alice", FilePath: "req.adoc", LineNumber: 5},
			"REQ-002": {ID: "REQ-002", ASPICE: "", Priority: "", ReviewedBy: "", FilePath: "req.adoc", LineNumber: 20},
		},
		ArchElements: map[string]*model.ArchElement{
			"comp.system": {ID: "comp.system", Parent: ""},
			"comp.api":    {ID: "comp.api", Parent: "comp.system", Req: []string{"REQ-001"}, Impl: "api.go", FilePath: "arc42.adoc", LineNumber: 10},
			"comp.db":     {ID: "comp.db", Parent: "comp.system", Req: []string{}, Impl: "", FilePath: "arc42.adoc", LineNumber: 20},
		},
		DesignElements: make(map[string]*model.DesignElement),
		TestSpecs: map[string]*model.TestSpec{
			"TS-001": {ID: "TS-001", Req: []string{"REQ-001"}, Arch: []string{"comp.api"}, FilePath: "test.adoc", LineNumber: 5},
		},
		TestCodes:   make(map[string]*model.TestCode),
		TestResults: make(map[string]*model.TestResult),
		Links: []*model.TraceLink{
			{FromID: "REQ-001", ToID: "comp.api", LinkType: "satisfied-by", Status: "active"},
			{FromID: "REQ-001", ToID: "TS-001", LinkType: "verified-by", Status: "active"},
			{FromID: "X", ToID: "Y", LinkType: "satisfied-by", Status: "stale"},
		},
	}
	return graph.NewAnalyzer(g)
}

func buildConfig(rules map[string]string) *model.Config {
	return &model.Config{Rules: rules, RuleParams: map[string]int{}}
}

// [test-spec,id=TS-VAL-001,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleEngine_Run verifies that Run processes all configured rules.
func TestRuleEngine_Run(t *testing.T) {
	cfg := buildConfig(map[string]string{
		"missing-review":            "warning",
		"all-reqs-must-have-aspice": "error",
	})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	if len(results) == 0 {
		t.Error("expected results, got none")
	}
}

// [test-spec,id=TS-VAL-002,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleEngine_Off verifies that rules set to 'off' produce no violations.
func TestRuleEngine_Off(t *testing.T) {
	cfg := buildConfig(map[string]string{
		"missing-review": "off",
	})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	for _, r := range results {
		if r.RuleID == "missing-review" && len(r.Violations) > 0 {
			t.Error("rule set to off should produce no violations")
		}
	}
}

// [test-spec,id=TS-VAL-003,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleMissingReview verifies that requirements without reviewed-by are flagged.
func TestRuleMissingReview(t *testing.T) {
	cfg := buildConfig(map[string]string{"missing-review": "warning"})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	var found bool
	for _, r := range results {
		for _, v := range r.Violations {
			if strings.Contains(v.Message, "REQ-002") {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected REQ-002 flagged for missing review")
	}
}

// [test-spec,id=TS-VAL-004,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleAllReqsMustHaveASPICE verifies that requirements without aspice are flagged.
func TestRuleAllReqsMustHaveASPICE(t *testing.T) {
	cfg := buildConfig(map[string]string{"all-reqs-must-have-aspice": "error"})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	numErrors, _ := TotalViolations(results)
	if numErrors == 0 {
		t.Error("expected error for REQ-002 missing aspice")
	}
}

// [test-spec,id=TS-VAL-005,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleMissingImpl verifies that arch elements without impl are flagged.
func TestRuleMissingImpl(t *testing.T) {
	cfg := buildConfig(map[string]string{"missing-impl": "warning"})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	var found bool
	for _, r := range results {
		for _, v := range r.Violations {
			if strings.Contains(v.Message, "comp.db") {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected comp.db flagged for missing impl")
	}
}

// [test-spec,id=TS-VAL-006,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleOrphanArchitecture verifies that arch elements without reqs are flagged.
func TestRuleOrphanArchitecture(t *testing.T) {
	cfg := buildConfig(map[string]string{"orphan-architecture": "warning"})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	var found bool
	for _, r := range results {
		for _, v := range r.Violations {
			if strings.Contains(v.Message, "comp.db") {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected comp.db flagged as orphan architecture")
	}
}

// [test-spec,id=TS-VAL-007,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleStaleTraces verifies that stale links are flagged.
func TestRuleStaleTraces(t *testing.T) {
	cfg := buildConfig(map[string]string{"stale-traces": "warning"})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	var found bool
	for _, r := range results {
		if r.RuleID == "stale-traces" && len(r.Violations) > 0 {
			found = true
		}
	}
	if !found {
		t.Error("expected stale trace violation")
	}
}

// [test-spec,id=TS-VAL-008,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleMaxOrphanPercentage verifies threshold-based orphan detection.
func TestRuleMaxOrphanPercentage(t *testing.T) {
	// REQ-002 is orphan → 50% orphan rate; threshold 30% → should trigger
	cfg := &model.Config{
		Rules:      map[string]string{"max-orphan-percentage": "warning"},
		RuleParams: map[string]int{"max-orphan-percentage": 30},
	}
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	var found bool
	for _, r := range results {
		if r.RuleID == "max-orphan-percentage" && len(r.Violations) > 0 {
			found = true
		}
	}
	if !found {
		t.Error("expected max-orphan-percentage violation at 50% > 30%")
	}
}

// [test-spec,id=TS-VAL-009,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleMaxOrphanPercentage_BelowThreshold verifies no violation when below threshold.
func TestRuleMaxOrphanPercentage_BelowThreshold(t *testing.T) {
	cfg := &model.Config{
		Rules:      map[string]string{"max-orphan-percentage": "warning"},
		RuleParams: map[string]int{"max-orphan-percentage": 60}, // 60% threshold, 50% actual
	}
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	for _, r := range results {
		if r.RuleID == "max-orphan-percentage" && len(r.Violations) > 0 {
			t.Error("expected no violation when below threshold")
		}
	}
}

// [test-spec,id=TS-VAL-010,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleUnknown verifies that unknown rule IDs produce a warning.
func TestRuleUnknown(t *testing.T) {
	cfg := buildConfig(map[string]string{"my-custom-unknown-rule": "error"})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	var found bool
	for _, r := range results {
		for _, v := range r.Violations {
			if strings.Contains(v.Message, "unknown rule") {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected unknown rule warning")
	}
}

// [test-spec,id=TS-VAL-011,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestTotalViolations verifies counting of errors and warnings.
func TestTotalViolations(t *testing.T) {
	results := []*RuleResult{
		{Violations: []Violation{
			{Severity: SeverityError},
			{Severity: SeverityWarning},
			{Severity: SeverityWarning},
		}},
		{Violations: []Violation{
			{Severity: SeverityError},
		}},
	}
	errs, warns := TotalViolations(results)
	if errs != 2 {
		t.Errorf("errors = %d, want 2", errs)
	}
	if warns != 2 {
		t.Errorf("warnings = %d, want 2", warns)
	}
}

// [test-spec,id=TS-VAL-012,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestFormatResults verifies that results are formatted with correct prefixes.
func TestFormatResults(t *testing.T) {
	results := []*RuleResult{{
		Violations: []Violation{
			{Severity: SeverityError, Message: "something broke", Location: "file.adoc:5"},
			{Severity: SeverityWarning, Message: "just a hint"},
		},
	}}
	out := FormatResults(results)
	if !strings.Contains(out, "[ERROR]") {
		t.Error("expected [ERROR] prefix")
	}
	if !strings.Contains(out, "[WARNING]") {
		t.Error("expected [WARNING] prefix")
	}
	if !strings.Contains(out, "file.adoc:5") {
		t.Error("expected location in output")
	}
}

// [test-spec,id=TS-VAL-013,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestParseSeverity verifies severity string parsing.
func TestParseSeverity(t *testing.T) {
	cases := []struct{ in string; want Severity }{
		{"error", SeverityError},
		{"ERROR", SeverityError},
		{"warning", SeverityWarning},
		{"off", SeverityOff},
		{"disabled", SeverityOff},
		{"unknown", SeverityWarning},
	}
	for _, tc := range cases {
		got := parseSeverity(tc.in)
		if got != tc.want {
			t.Errorf("parseSeverity(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// [test-spec,id=TS-VAL-014,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleMissingTestSpec verifies that arch elements without test links are flagged.
func TestRuleMissingTestSpec(t *testing.T) {
	cfg := buildConfig(map[string]string{"missing-test-spec": "warning"})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	var found bool
	for _, r := range results {
		for _, v := range r.Violations {
			if strings.Contains(v.Message, "comp.db") {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected comp.db flagged for missing test spec")
	}
}

// [test-spec,id=TS-VAL-015,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleEngine_NilRulesMap verifies Run() is safe when config.Rules is nil.
func TestRuleEngine_NilRulesMap(t *testing.T) {
	cfg := &model.Config{Rules: nil, RuleParams: map[string]int{}}
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run() // must not panic
	if len(results) != 0 {
		t.Errorf("expected 0 results for nil rules map, got %d", len(results))
	}
}

// [test-spec,id=TS-VAL-016,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleMaxOrphanPercentage_NilRuleParams verifies default threshold when RuleParams is nil.
func TestRuleMaxOrphanPercentage_NilRuleParams(t *testing.T) {
	cfg := &model.Config{
		Rules:      map[string]string{"max-orphan-percentage": "warning"},
		RuleParams: nil, // nil map — should use default threshold of 20%
	}
	// REQ-002 is orphan → 50% orphan rate, default threshold 20% → should trigger
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	var found bool
	for _, r := range results {
		if r.RuleID == "max-orphan-percentage" && len(r.Violations) > 0 {
			found = true
		}
	}
	if !found {
		t.Error("expected max-orphan-percentage violation with nil RuleParams (uses default 20%)")
	}
}

// [test-spec,id=TS-VAL-017,req="REQ-VALIDATE-001",aspice="SWE.5.BP3"]
// TestRuleUnknownBausteinsichtRule verifies undocumented-bausteinsicht-elements produces "unknown rule" warning.
func TestRuleUnknownBausteinsichtRule(t *testing.T) {
	cfg := buildConfig(map[string]string{"undocumented-bausteinsicht-elements": "error"})
	engine := NewRuleEngine(cfg, buildValidationGraph())
	results := engine.Run()
	var found bool
	for _, r := range results {
		for _, v := range r.Violations {
			if strings.Contains(v.Message, "unknown rule") {
				found = true
			}
		}
	}
	if !found {
		t.Error("expected 'unknown rule' warning for undocumented-bausteinsicht-elements")
	}
}
