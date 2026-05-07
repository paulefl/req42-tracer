package report

import (
	"testing"

	"github.com/paulefl/req42-tracer/internal/model"
)

func TestExportGraphData(t *testing.T) {
	// Create a minimal graph
	graph := &model.TraceabilityGraph{
		Requirements: map[string]*model.Requirement{
			"REQ-001": {
				ID:       "REQ-001",
				Title:    "Test Requirement",
				Priority: "high",
				Status:   "approved",
				ASPICE:   "SWE.1",
			},
		},
		ArchElements: map[string]*model.ArchElement{
			"comp.api": {
				ID:    "comp.api",
				Title: "API Component",
				ASPICE: "SWE.2",
				Req:    []string{"REQ-001"},
			},
		},
		TestSpecs: map[string]*model.TestSpec{
			"spec.api": {
				ID:    "spec.api",
				Title: "API Test Spec",
				Req:   []string{"REQ-001"},
				Arch:  []string{"comp.api"},
			},
		},
		TestResults: map[string]*model.TestResult{
			"result-1": {
				ID:       "result-1",
				TestName: "TestAPI",
				Status:   "passed",
			},
		},
		Links: []*model.TraceLink{
			{
				FromID:   "REQ-001",
				FromType: "requirement",
				ToID:     "comp.api",
				ToType:   "arch",
				LinkType: "satisfies",
			},
		},
	}

	// Export graph data
	data := ExportGraphData(graph)

	// Validate
	if len(data.Nodes) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(data.Nodes))
	}

	if len(data.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(data.Edges))
	}

	// Check node types
	nodeTypes := make(map[string]int)
	for _, node := range data.Nodes {
		nodeTypes[node.Type]++
	}

	if nodeTypes["requirement"] != 1 {
		t.Errorf("Expected 1 requirement node, got %d", nodeTypes["requirement"])
	}

	if nodeTypes["arch"] != 1 {
		t.Errorf("Expected 1 arch node, got %d", nodeTypes["arch"])
	}

	if nodeTypes["test-spec"] != 1 {
		t.Errorf("Expected 1 test-spec node, got %d", nodeTypes["test-spec"])
	}

	if nodeTypes["test-result"] != 1 {
		t.Errorf("Expected 1 test-result node, got %d", nodeTypes["test-result"])
	}

	// Check edge
	edge := data.Edges[0]
	if edge.Label != "satisfies" {
		t.Errorf("Expected edge label 'satisfies', got %s", edge.Label)
	}

	t.Log("✅ All graph export tests passed")
}

func TestFilterGraphByType(t *testing.T) {
	// Create graph data
	data := &GraphData{
		Nodes: []Node{
			{ID: "REQ-001", Type: "requirement"},
			{ID: "arch-001", Type: "arch"},
			{ID: "spec-001", Type: "test-spec"},
		},
		Edges: []Edge{
			{Source: "REQ-001", Target: "arch-001", Label: "satisfies"},
			{Source: "arch-001", Target: "spec-001", Label: "implements"},
		},
	}

	// Filter by requirement only
	filtered := FilterGraphByType(data, "requirement")

	if len(filtered.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(filtered.Nodes))
	}

	if len(filtered.Edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(filtered.Edges))
	}

	// Filter by requirement and arch
	filtered = FilterGraphByType(data, "requirement", "arch")

	if len(filtered.Nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(filtered.Nodes))
	}

	if len(filtered.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(filtered.Edges))
	}

	t.Log("✅ All filter tests passed")
}
