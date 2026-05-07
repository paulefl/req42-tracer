package report

import (
	"fmt"

	"github.com/paulefl/req42-tracer/internal/model"
)

// Node represents a graph node for D3.js visualization.
type Node struct {
	ID       string                 `json:"id"`
	Label    string                 `json:"label"`
	Type     string                 `json:"type"` // "requirement", "arch", "test-spec", "test-code", "test-result"
	Metadata map[string]interface{} `json:"metadata"`
	Group    int                    `json:"group"` // For D3 force layout grouping
}

// Edge represents a graph edge for D3.js visualization.
type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"` // satisfies, implements, verifies, derives, covers
	Value  int    `json:"value"` // Weight for D3
}

// GraphData represents the complete graph for D3.js.
type GraphData struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// ExportGraphData converts a TraceabilityGraph to D3-compatible GraphData.
func ExportGraphData(g *model.TraceabilityGraph) *GraphData {
	data := &GraphData{
		Nodes: make([]Node, 0),
		Edges: make([]Edge, 0),
	}

	nodeMap := make(map[string]bool) // Track added nodes to avoid duplicates

	// Add requirement nodes (group 0)
	for _, req := range g.Requirements {
		node := Node{
			ID:    req.ID,
			Label: fmt.Sprintf("%s: %s", req.ID, req.Title),
			Type:  "requirement",
			Group: 0,
			Metadata: map[string]interface{}{
				"priority": req.Priority,
				"status":   req.Status,
				"aspice":   req.ASPICE,
				"version":  req.Version,
			},
		}
		data.Nodes = append(data.Nodes, node)
		nodeMap[req.ID] = true
	}

	// Add architecture nodes (group 1)
	for _, arch := range g.ArchElements {
		node := Node{
			ID:    arch.ID,
			Label: fmt.Sprintf("%s: %s", arch.ID, arch.Title),
			Type:  "arch",
			Group: 1,
			Metadata: map[string]interface{}{
				"aspice": arch.ASPICE,
				"impl":   arch.Impl,
				"parent": arch.Parent,
			},
		}
		data.Nodes = append(data.Nodes, node)
		nodeMap[arch.ID] = true
	}

	// Add test specification nodes (group 2)
	for _, spec := range g.TestSpecs {
		node := Node{
			ID:    spec.ID,
			Label: fmt.Sprintf("%s: %s", spec.ID, spec.Title),
			Type:  "test-spec",
			Group: 2,
			Metadata: map[string]interface{}{
				"req":  spec.Req,
				"arch": spec.Arch,
			},
		}
		data.Nodes = append(data.Nodes, node)
		nodeMap[spec.ID] = true
	}

	// Add test code nodes (group 2)
	for _, code := range g.TestCodes {
		node := Node{
			ID:    code.ID,
			Label: fmt.Sprintf("%s: %s", code.ID, code.Function),
			Type:  "test-code",
			Group: 2,
			Metadata: map[string]interface{}{
				"test_spec": code.TestSpec,
				"file":      code.File,
				"language":  code.Language,
			},
		}
		data.Nodes = append(data.Nodes, node)
		nodeMap[code.ID] = true
	}

	// Add test result nodes (group 3)
	for _, result := range g.TestResults {
		node := Node{
			ID:    result.ID,
			Label: fmt.Sprintf("%s: %s (%s)", result.ID, result.TestName, result.Status),
			Type:  "test-result",
			Group: 3,
			Metadata: map[string]interface{}{
				"status":       result.Status,
				"duration":     result.Duration,
				"linked_spec":  result.LinkedSpec,
				"linked_code":  result.LinkedCode,
				"platform":     result.Platform,
				"pass_rate":    0, // Will be set later
			},
		}
		data.Nodes = append(data.Nodes, node)
		nodeMap[result.ID] = true
	}

	// Add edges from trace links
	linkTypeWeight := map[string]int{
		"satisfies":  2,
		"implements": 2,
		"verifies":   3,
		"derives":    1,
		"covers":     2,
	}

	for _, link := range g.Links {
		// Only add edges for nodes that exist
		if nodeMap[link.FromID] && nodeMap[link.ToID] {
			weight := 1
			if w, ok := linkTypeWeight[link.LinkType]; ok {
				weight = w
			}

			edge := Edge{
				Source: link.FromID,
				Target: link.ToID,
				Label:  link.LinkType,
				Value:  weight,
			}
			data.Edges = append(data.Edges, edge)
		}
	}

	return data
}

// FilterGraphByType returns a new GraphData with only specified node types.
func FilterGraphByType(data *GraphData, types ...string) *GraphData {
	typeMap := make(map[string]bool)
	for _, t := range types {
		typeMap[t] = true
	}

	filtered := &GraphData{
		Nodes: make([]Node, 0),
		Edges: make([]Edge, 0),
	}

	nodeMap := make(map[string]bool)
	for _, node := range data.Nodes {
		if typeMap[node.Type] {
			filtered.Nodes = append(filtered.Nodes, node)
			nodeMap[node.ID] = true
		}
	}

	// Only add edges where both source and target are in filtered nodes
	for _, edge := range data.Edges {
		if nodeMap[edge.Source] && nodeMap[edge.Target] {
			filtered.Edges = append(filtered.Edges, edge)
		}
	}

	return filtered
}
