package graph

import (
	"fmt"
	"strings"

	"github.com/paulefl/req42-tracer/internal/model"
)

// Builder constructs the traceability graph from parsed artifacts.
type Builder struct {
	graph *model.TraceabilityGraph
}

// NewBuilder creates a new graph builder.
func NewBuilder() *Builder {
	return &Builder{
		graph: &model.TraceabilityGraph{
			Requirements: make(map[string]*model.Requirement),
			ArchElements: make(map[string]*model.ArchElement),
			TestSpecs:    make(map[string]*model.TestSpec),
			TestCodes:    make(map[string]*model.TestCode),
			TestResults:  make(map[string]*model.TestResult),
			Links:        []*model.TraceLink{},
		},
	}
}

// MergeGraph merges another graph into this one.
func (b *Builder) MergeGraph(other *model.TraceabilityGraph) error {
	if other == nil {
		return nil
	}

	// Merge all artifact types
	for id, req := range other.Requirements {
		if _, exists := b.graph.Requirements[id]; exists {
			return fmt.Errorf("duplicate requirement ID: %s", id)
		}
		b.graph.Requirements[id] = req
	}

	for id, arch := range other.ArchElements {
		if _, exists := b.graph.ArchElements[id]; exists {
			return fmt.Errorf("duplicate architecture element ID: %s", id)
		}
		b.graph.ArchElements[id] = arch
	}

	for id, spec := range other.TestSpecs {
		if _, exists := b.graph.TestSpecs[id]; exists {
			return fmt.Errorf("duplicate test-spec ID: %s", id)
		}
		b.graph.TestSpecs[id] = spec
	}

	for id, code := range other.TestCodes {
		if _, exists := b.graph.TestCodes[id]; exists {
			return fmt.Errorf("duplicate test-code ID: %s", id)
		}
		b.graph.TestCodes[id] = code
	}

	for id, result := range other.TestResults {
		if _, exists := b.graph.TestResults[id]; exists {
			return fmt.Errorf("duplicate test-result ID: %s", id)
		}
		b.graph.TestResults[id] = result
	}

	b.graph.Links = append(b.graph.Links, other.Links...)

	return nil
}

// BuildLinks constructs trace links between artifacts based on explicit attributes and references.
func (b *Builder) BuildLinks() error {
	// Link requirements to architecture
	for archID, arch := range b.graph.ArchElements {
		for _, reqID := range arch.Req {
			if _, exists := b.graph.Requirements[reqID]; !exists {
				continue // Skip non-existent requirements (will be flagged as error elsewhere)
			}

			b.addLink(&model.TraceLink{
				FromID:   reqID,
				FromType: "requirement",
				ToID:     archID,
				ToType:   "arch",
				LinkType: "satisfied-by",
				Status:   "active",
				Reason:   "Explicit ref= attribute",
			})
		}
	}

	// Link architecture to test specs
	for specID, spec := range b.graph.TestSpecs {
		for _, archID := range spec.Arch {
			if _, exists := b.graph.ArchElements[archID]; !exists {
				continue
			}

			b.addLink(&model.TraceLink{
				FromID:   archID,
				FromType: "arch",
				ToID:     specID,
				ToType:   "test-spec",
				LinkType: "verified-by",
				Status:   "active",
				Reason:   "Explicit arch= attribute",
			})
		}
	}

	// Link test specs to requirements
	for specID, spec := range b.graph.TestSpecs {
		for _, reqID := range spec.Req {
			if _, exists := b.graph.Requirements[reqID]; !exists {
				continue
			}

			b.addLink(&model.TraceLink{
				FromID:   reqID,
				FromType: "requirement",
				ToID:     specID,
				ToType:   "test-spec",
				LinkType: "verified-by",
				Status:   "active",
				Reason:   "Explicit req= attribute",
			})
		}
	}

	// Link test codes to test specs
	for codeID, code := range b.graph.TestCodes {
		if code.TestSpec != "" {
			if _, exists := b.graph.TestSpecs[code.TestSpec]; !exists {
				continue
			}

			b.addLink(&model.TraceLink{
				FromID:   code.TestSpec,
				FromType: "test-spec",
				ToID:     codeID,
				ToType:   "test-code",
				LinkType: "implements",
				Status:   "active",
				Reason:   "Explicit test-spec= attribute",
			})
		}
	}

	// Link test results to test codes (via name matching as fallback)
	for resultID, result := range b.graph.TestResults {
		if result.LinkedCode == "" && result.LinkedSpec == "" {
			// Try name-based matching: test name -> TestCode function name
			for codeID, code := range b.graph.TestCodes {
				if nameMatches(result.TestName, code.Function) {
					result.LinkedCode = codeID
					b.addLink(&model.TraceLink{
						FromID:   codeID,
						FromType: "test-code",
						ToID:     resultID,
						ToType:   "test-result",
						LinkType: "produces",
						Status:   "active",
						Reason:   "Name-based matching",
					})
					break
				}
			}
		}

		// If still unlinked, try spec matching
		if result.LinkedSpec == "" {
			for specID := range b.graph.TestSpecs {
				// Check if test name matches spec ID pattern
				if strings.Contains(strings.ToLower(result.TestName), strings.ToLower(specID)) {
					result.LinkedSpec = specID
					b.addLink(&model.TraceLink{
						FromID:   specID,
						FromType: "test-spec",
						ToID:     resultID,
						ToType:   "test-result",
						LinkType: "produces",
						Status:   "active",
						Reason:   "Name-based spec matching",
					})
					break
				}
			}
		}
	}

	return nil
}

// addLink adds a trace link, avoiding duplicates.
func (b *Builder) addLink(link *model.TraceLink) {
	// Check if link already exists
	for _, existing := range b.graph.Links {
		if existing.FromID == link.FromID && existing.ToID == link.ToID &&
			existing.FromType == link.FromType && existing.ToType == link.ToType {
			return // Link already exists
		}
	}
	b.graph.Links = append(b.graph.Links, link)
}

// DeriveASPICELevels automatically derives ASPICE levels from ARC42 hierarchy.
func (b *Builder) DeriveASPICELevels() {
	for _, arch := range b.graph.ArchElements {
		if arch.ASPICE != "" {
			continue // Already explicitly set
		}

		if arch.Parent == "" {
			// Level 1 (no parent) -> SWE.2
			arch.ASPICE = "SWE.2"
		} else {
			// Level 2+ (has parent) -> SWE.3
			arch.ASPICE = "SWE.3"
		}
	}
}

// GetGraph returns the built traceability graph.
func (b *Builder) GetGraph() *model.TraceabilityGraph {
	return b.graph
}

// nameMatches checks if a test result name matches a test code function name.
// Example: "TestAPIAuth" matches function "TestAPIAuth" or "test_api_auth"
func nameMatches(resultName, funcName string) bool {
	// Normalize both to lowercase for comparison
	rLower := strings.ToLower(resultName)
	fLower := strings.ToLower(funcName)

	// Exact match
	if rLower == fLower {
		return true
	}

	// Remove common prefixes/suffixes
	rNorm := strings.TrimPrefix(rLower, "test")
	rNorm = strings.TrimPrefix(rNorm, "test_")
	fNorm := strings.TrimPrefix(fLower, "test")
	fNorm = strings.TrimPrefix(fNorm, "test_")

	// Convert underscores to empty for comparison
	rNorm = strings.ReplaceAll(rNorm, "_", "")
	fNorm = strings.ReplaceAll(fNorm, "_", "")

	return rNorm != "" && rNorm == fNorm
}
