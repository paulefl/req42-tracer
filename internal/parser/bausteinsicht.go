package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/paulefl/req42-tracer/internal/model"
)

// BausteinsichtParser loads and parses Bausteinsicht architecture models (JSONC).
type BausteinsichtParser struct {
	filePath string
}

// NewBausteinsichtParser creates a new parser for a Bausteinsicht model file.
func NewBausteinsichtParser(filePath string) *BausteinsichtParser {
	return &BausteinsichtParser{filePath: filePath}
}

// BausteinsichtModel represents the structure of an architecture.jsonc file.
type BausteinsichtModel struct {
	Model map[string]*BausteinsichtProject `json:"model"`
	Views []*BausteinsichtView              `json:"views,omitempty"`
}

// BausteinsichtProject represents a project in the architecture model.
type BausteinsichtProject struct {
	Description string                       `json:"description"`
	Elements    map[string]*BausteinsichtElement `json:"elements"`
}

// BausteinsichtElement represents an element in the architecture.
type BausteinsichtElement struct {
	Description string                       `json:"description"`
	Type        string                       `json:"type,omitempty"`
	Technology  string                       `json:"technology,omitempty"`
	Parent      string                       `json:"parent,omitempty"`
	Elements    map[string]*BausteinsichtElement `json:"elements,omitempty"`
}

// BausteinsichtView represents a diagram view in the architecture.
type BausteinsichtView struct {
	Name    string   `json:"name"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude,omitempty"`
}

// Parse loads and parses the Bausteinsicht model file.
func (p *BausteinsichtParser) Parse(project string) (*model.TraceabilityGraph, error) {
	data, err := os.ReadFile(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Bausteinsicht model %s: %w", p.filePath, err)
	}

	// Remove JSONC comments before parsing
	jsonStr := removeComments(string(data))

	var bausteinsichtModel BausteinsichtModel
	if err := json.Unmarshal([]byte(jsonStr), &bausteinsichtModel); err != nil {
		return nil, fmt.Errorf("failed to parse Bausteinsicht model %s: %w", p.filePath, err)
	}

	graph := &model.TraceabilityGraph{
		Requirements: make(map[string]*model.Requirement),
		ArchElements: make(map[string]*model.ArchElement),
		TestSpecs:    make(map[string]*model.TestSpec),
		TestCodes:    make(map[string]*model.TestCode),
		TestResults:  make(map[string]*model.TestResult),
		Links:        []*model.TraceLink{},
	}

	// Process all projects and their elements
	for projName, projModel := range bausteinsichtModel.Model {
		p.flattenElements(projName, projModel.Elements, "", project, graph)
	}

	return graph, nil
}

// flattenElements recursively flattens the nested element hierarchy into dot-notation IDs.
func (p *BausteinsichtParser) flattenElements(
	parentPath string,
	elements map[string]*BausteinsichtElement,
	prefix string,
	project string,
	graph *model.TraceabilityGraph,
) {
	for elemName, elem := range elements {
		// Build dot-notation ID: system, backend, backend.parser, etc.
		var elemID string
		if prefix == "" {
			elemID = elemName
		} else {
			elemID = prefix + "." + elemName
		}

		// Create architecture element
		archElem := &model.ArchElement{
			ID:         elemID,
			Title:      elemName,
			Text:       elem.Description,
			Project:    project,
			FilePath:   p.filePath,
			Attributes: make(map[string]string),
		}

		if elem.Type != "" {
			archElem.Attributes["type"] = elem.Type
		}
		if elem.Technology != "" {
			archElem.Attributes["technology"] = elem.Technology
		}
		if elem.Parent != "" {
			archElem.Parent = elem.Parent
		}

		graph.ArchElements[elemID] = archElem

		// Recursively process nested elements
		if len(elem.Elements) > 0 {
			p.flattenElements(parentPath, elem.Elements, elemID, project, graph)
		}
	}
}

// removeComments removes JSONC comments from a string.
func removeComments(jsonc string) string {
	var result strings.Builder
	inString := false
	escapeNext := false

	lines := strings.Split(jsonc, "\n")
	for _, line := range lines {
		for i, char := range line {
			if escapeNext {
				result.WriteRune(char)
				escapeNext = false
				continue
			}

			if char == '\\' && inString {
				result.WriteRune(char)
				escapeNext = true
				continue
			}

			if char == '"' && !escapeNext {
				inString = !inString
				result.WriteRune(char)
				continue
			}

			if !inString && i+1 < len(line) {
				if char == '/' && line[i+1] == '/' {
					// Line comment; skip rest of line
					break
				}
				if char == '/' && line[i+1] == '*' {
					// Block comment; skip until */
					// For simplicity, assume block comments don't nest
					for j := i + 2; j < len(line); j++ {
						if j+1 < len(line) && line[j] == '*' && line[j+1] == '/' {
							// Skip to after */
							i = j + 1
							break
						}
					}
					continue
				}
			}

			result.WriteRune(char)
		}

		if !inString {
			result.WriteString("\n")
		}
	}

	return result.String()
}
