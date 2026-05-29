package lsp

import (
	"path/filepath"
	"strings"

	"github.com/paulefl/req42-tracer/src/internal/model"
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
	// lineNumber == 0 is the Go zero-value (unset); model stores 1-based numbers.
	// Return nil so the client doesn't silently jump to line 0 of the wrong file.
	if lineNumber == 0 {
		return nil
	}

	uri := pathToURI(filePath)
	line := lineNumber - 1 // convert 1-based → 0-based

	return &Location{
		URI: uri,
		Range: diagRange{
			Start: diagPosition{Line: line, Character: 0},
			End:   diagPosition{Line: line, Character: 0},
		},
	}
}

// pathToURI converts a file system path to a file:// URI per RFC 8089.
// Unix:    /abs/path       → file:///abs/path
// Windows: C:\abs\path     → file:///C:/abs/path  (empty authority + drive)
func pathToURI(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	abs = filepath.ToSlash(abs)
	if strings.HasPrefix(abs, "/") {
		return "file://" + abs // Unix: leading '/' gives file:///abs/path
	}
	return "file:///" + abs // Windows: prepend empty authority for file:///C:/...
}
