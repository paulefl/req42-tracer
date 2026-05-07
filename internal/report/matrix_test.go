package report

import (
	"testing"

	"github.com/paulefl/req42-tracer/internal/model"
)

func TestBuildMatrixData(t *testing.T) {
	// Create a test graph
	graph := &model.TraceabilityGraph{
		Requirements: map[string]*model.Requirement{
			"REQ-001": {ID: "REQ-001", Title: "Test Req 1", Priority: "high", Status: "approved"},
			"REQ-002": {ID: "REQ-002", Title: "Test Req 2", Priority: "low", Status: "draft"},
		},
		ArchElements: map[string]*model.ArchElement{
			"arch-001": {ID: "arch-001", Title: "Arch 1"},
			"arch-002": {ID: "arch-002", Title: "Arch 2"},
		},
		TestSpecs: map[string]*model.TestSpec{
			"spec-001": {ID: "spec-001", Title: "Spec 1"},
		},
		TestResults:  make(map[string]*model.TestResult),
		TestCodes:    make(map[string]*model.TestCode),
		Links: []*model.TraceLink{
			{FromID: "REQ-001", FromType: "requirement", ToID: "arch-001", ToType: "arch", LinkType: "satisfies", Status: "active"},
			{FromID: "REQ-001", FromType: "requirement", ToID: "spec-001", ToType: "test-spec", LinkType: "verifies", Status: "active"},
		},
	}

	matrix := BuildMatrixData(graph)

	if len(matrix.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(matrix.Rows))
	}

	if len(matrix.Columns) != 3 {
		t.Errorf("Expected 3 columns (2 arch + 1 spec), got %d", len(matrix.Columns))
	}

	if matrix.Statistics.TotalRequirements != 2 {
		t.Errorf("Expected 2 total requirements, got %d", matrix.Statistics.TotalRequirements)
	}

	if matrix.Statistics.CoveredRequirements != 1 {
		t.Errorf("Expected 1 covered requirement, got %d", matrix.Statistics.CoveredRequirements)
	}

	t.Log("✅ Matrix data structure test passed")
}

func TestFilterMatrix(t *testing.T) {
	graph := &model.TraceabilityGraph{
		Requirements: map[string]*model.Requirement{
			"REQ-001": {ID: "REQ-001", Title: "High Req", Priority: "high", Status: "approved"},
			"REQ-002": {ID: "REQ-002", Title: "Low Req", Priority: "low", Status: "draft"},
		},
		ArchElements: map[string]*model.ArchElement{
			"arch-001": {ID: "arch-001", Title: "Arch 1"},
		},
		TestSpecs:   make(map[string]*model.TestSpec),
		TestResults: make(map[string]*model.TestResult),
		TestCodes:   make(map[string]*model.TestCode),
		Links:       []*model.TraceLink{},
	}

	matrix := BuildMatrixData(graph)

	// Filter by high priority
	filtered := FilterMatrix(matrix, []string{"high"}, []string{"approved", "draft"})

	if len(filtered.Rows) != 1 {
		t.Errorf("Expected 1 filtered row, got %d", len(filtered.Rows))
	}

	if filtered.Rows[0].RequirementID != "REQ-001" {
		t.Errorf("Expected REQ-001, got %s", filtered.Rows[0].RequirementID)
	}

	t.Log("✅ Matrix filtering test passed")
}

func TestExportMatrixCSV(t *testing.T) {
	graph := &model.TraceabilityGraph{
		Requirements: map[string]*model.Requirement{
			"REQ-001": {ID: "REQ-001", Title: "Test", Priority: "high", Status: "approved"},
		},
		ArchElements: map[string]*model.ArchElement{
			"arch-001": {ID: "arch-001", Title: "Arch"},
		},
		TestSpecs:   make(map[string]*model.TestSpec),
		TestResults: make(map[string]*model.TestResult),
		TestCodes:   make(map[string]*model.TestCode),
		Links: []*model.TraceLink{
			{FromID: "REQ-001", FromType: "requirement", ToID: "arch-001", ToType: "arch", LinkType: "satisfies"},
		},
	}

	matrix := BuildMatrixData(graph)
	csv := ExportMatrixToCSV(matrix)

	if !Contains(csv, "Requirement,Priority,Status") {
		t.Error("CSV header not found")
	}

	if !Contains(csv, "REQ-001,high,approved") {
		t.Error("CSV requirement row not found")
	}

	t.Log("✅ CSV export test passed")
}

func Contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
