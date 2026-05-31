package report

import (
	"testing"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

// [test-spec,id=TS-RPT-046,req="REQ-REPORT-001",aspice="SWE.5.BP3"]
// TestBuildElementsData_Types verifies all element types appear in the output.
func TestBuildElementsData_Types(t *testing.T) {
	g := &model.TraceabilityGraph{
		Requirements: map[string]*model.Requirement{
			"REQ-001": {ID: "REQ-001", Title: "Req One"},
		},
		ArchElements: map[string]*model.ArchElement{
			"comp.a": {ID: "comp.a", Title: "Comp A", Impl: "internal/a"},
		},
		DesignElements: map[string]*model.DesignElement{
			"DSN-001": {ID: "DSN-001", Title: "Design One"},
		},
		TestSpecs: map[string]*model.TestSpec{
			"TS-001": {ID: "TS-001", Title: "Spec One"},
		},
		TestResults: map[string]*model.TestResult{
			"a::TestFoo": {ID: "a::TestFoo", TestName: "TestFoo", Status: "passed"},
		},
		TestCodes: make(map[string]*model.TestCode),
		Links:     []*model.TraceLink{},
	}

	data := BuildElementsData(g)

	types := make(map[string]int)
	for _, item := range data.Items {
		types[item.Type]++
	}
	for _, want := range []string{"req", "arch", "dsn", "test-spec", "test-result"} {
		if types[want] == 0 {
			t.Errorf("expected element of type %q, not found", want)
		}
	}
}

// [test-spec,id=TS-RPT-047,req="REQ-REPORT-001",aspice="SWE.5.BP3"]
// TestBuildElementsData_TraceLinks verifies trace_up and trace_down are populated from Links.
func TestBuildElementsData_TraceLinks(t *testing.T) {
	g := &model.TraceabilityGraph{
		Requirements: map[string]*model.Requirement{
			"REQ-001": {ID: "REQ-001", Title: "Req"},
		},
		ArchElements: map[string]*model.ArchElement{
			"comp.a": {ID: "comp.a", Title: "Comp"},
		},
		DesignElements: make(map[string]*model.DesignElement),
		TestSpecs:      make(map[string]*model.TestSpec),
		TestResults:    make(map[string]*model.TestResult),
		TestCodes:      make(map[string]*model.TestCode),
		Links: []*model.TraceLink{
			{FromID: "REQ-001", FromType: "requirement", ToID: "comp.a", ToType: "arch", LinkType: "satisfies"},
		},
	}

	data := BuildElementsData(g)

	var req, arch *ElementItem
	for i := range data.Items {
		switch data.Items[i].ID {
		case "REQ-001":
			req = &data.Items[i]
		case "comp.a":
			arch = &data.Items[i]
		}
	}
	if req == nil || arch == nil {
		t.Fatal("expected both REQ-001 and comp.a in elements")
	}
	if len(req.TraceDown) != 1 || req.TraceDown[0] != "comp.a" {
		t.Errorf("REQ-001 trace_down = %v, want [comp.a]", req.TraceDown)
	}
	if len(arch.TraceUp) != 1 || arch.TraceUp[0] != "REQ-001" {
		t.Errorf("comp.a trace_up = %v, want [REQ-001]", arch.TraceUp)
	}
}

// [test-spec,id=TS-RPT-048,req="REQ-REPORT-001",aspice="SWE.5.BP3"]
// TestBuildElementsData_EmptyGraph verifies an empty graph produces an empty items list.
func TestBuildElementsData_EmptyGraph(t *testing.T) {
	g := &model.TraceabilityGraph{
		Requirements:   make(map[string]*model.Requirement),
		ArchElements:   make(map[string]*model.ArchElement),
		DesignElements: make(map[string]*model.DesignElement),
		TestSpecs:      make(map[string]*model.TestSpec),
		TestResults:    make(map[string]*model.TestResult),
		TestCodes:      make(map[string]*model.TestCode),
		Links:          []*model.TraceLink{},
	}
	data := BuildElementsData(g)
	if len(data.Items) != 0 {
		t.Errorf("expected 0 items for empty graph, got %d", len(data.Items))
	}
}
