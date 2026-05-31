package report

import (
	"testing"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

// [test-spec,id=TS-RPT-032,req="REQ-REPORT-002",aspice="SWE.5.BP3"]
// TestBuildGapsData_Empty verifies that an empty GapAnalysisResult produces no gaps.
func TestBuildGapsData_Empty(t *testing.T) {
	gaps := &model.GapAnalysisResult{
		OrphanRequirements:     []*model.Requirement{},
		OrphanArchElements:     []*model.ArchElement{},
		OrphanDesignElements:   []*model.DesignElement{},
		UntestedArchElements:   []*model.ArchElement{},
		UntestedDesignElements: []*model.DesignElement{},
		OrphanTestSpecs:        []*model.TestSpec{},
		UntracedTestResults:    []*model.TestResult{},
	}

	data := BuildGapsData(gaps)

	if data.HasGaps {
		t.Error("HasGaps = true, want false for empty result")
	}
	if len(data.OrphanRequirements) != 0 {
		t.Errorf("OrphanRequirements = %d, want 0", len(data.OrphanRequirements))
	}
}

// [test-spec,id=TS-RPT-033,req="REQ-REPORT-002",aspice="SWE.5.BP3"]
// TestBuildGapsData_WithGaps verifies that gaps are correctly mapped to GapItems.
func TestBuildGapsData_WithGaps(t *testing.T) {
	gaps := &model.GapAnalysisResult{
		OrphanRequirements: []*model.Requirement{
			{ID: "REQ-001", Title: "Orphan Req", Priority: "high"},
		},
		OrphanArchElements:     []*model.ArchElement{},
		OrphanDesignElements:   []*model.DesignElement{},
		UntestedArchElements:   []*model.ArchElement{},
		UntestedDesignElements: []*model.DesignElement{
			{ID: "DSN-001", Title: "Untested Design", Arch: "comp.parser"},
		},
		OrphanTestSpecs:     []*model.TestSpec{},
		UntracedTestResults: []*model.TestResult{},
	}

	data := BuildGapsData(gaps)

	if !data.HasGaps {
		t.Error("HasGaps = false, want true")
	}
	if len(data.OrphanRequirements) != 1 {
		t.Fatalf("OrphanRequirements = %d, want 1", len(data.OrphanRequirements))
	}
	if data.OrphanRequirements[0].ID != "REQ-001" {
		t.Errorf("ID = %q, want REQ-001", data.OrphanRequirements[0].ID)
	}
	if data.OrphanRequirements[0].Info != "high" {
		t.Errorf("Info = %q, want high", data.OrphanRequirements[0].Info)
	}
	if len(data.UntestedDesignElements) != 1 {
		t.Fatalf("UntestedDesignElements = %d, want 1", len(data.UntestedDesignElements))
	}
	if data.UntestedDesignElements[0].Info != "comp.parser" {
		t.Errorf("Info = %q, want comp.parser", data.UntestedDesignElements[0].Info)
	}
}
