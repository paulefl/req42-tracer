package lsp

import (
	"fmt"
	"path/filepath"

	"github.com/paulefl/req42-tracer/internal/model"
)

// Location is an LSP Location object (file + range).
type Location struct {
	URI   string    `json:"uri"`
	Range diagRange `json:"range"`
}

// findDefinition looks up attr+id in the graph and returns the LSP Location
// of its definition (the source .adoc file and line), or nil if not found.
func findDefinition(attr, id string, g *model.TraceabilityGraph) *Location {
	var filePath string
	var lineNumber int

	switch attr {
	case "req":
		r, ok := g.Requirements[id]
		if !ok {
			return nil
		}
		filePath, lineNumber = r.FilePath, r.LineNumber
	case "arch":
		a, ok := g.ArchElements[id]
		if !ok {
			return nil
		}
		filePath, lineNumber = a.FilePath, a.LineNumber
	case "test-spec":
		s, ok := g.TestSpecs[id]
		if !ok {
			return nil
		}
		filePath, lineNumber = s.FilePath, s.LineNumber
	default:
		return nil
	}

	if filePath == "" {
		return nil
	}

	// Convert OS path to LSP file URI.
	uri := pathToURI(filePath)

	// LSP line numbers are 0-based; model stores 1-based.
	line := lineNumber - 1
	if line < 0 {
		line = 0
	}

	return &Location{
		URI: uri,
		Range: diagRange{
			Start: diagPosition{Line: line, Character: 0},
			End:   diagPosition{Line: line, Character: 0},
		},
	}
}

// pathToURI converts a file system path to a file:// URI.
func pathToURI(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	// On all supported platforms the path separator is already '/'; on Windows
	// filepath.Abs uses '\' but LSP URIs require '/'.
	abs = filepath.ToSlash(abs)
	return fmt.Sprintf("file://%s", abs)
}
