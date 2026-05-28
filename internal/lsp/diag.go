package lsp

import (
	"regexp"

	"github.com/paulefl/req42-tracer/internal/model"
)

// diagPattern finds all req=/arch=/test-spec= attribute values in a line.
// Same structure as hoverPattern — group 1 = attr name, group 2 = value.
var diagPattern = regexp.MustCompile(`(?:^|[,\s])(req|arch|test-spec)=([^,\]"'\s]+)`)

// Diagnostic severity levels (LSP spec §3.17.1).
const diagError = 1

// Diagnostic is a single LSP Diagnostic object.
type Diagnostic struct {
	Range    diagRange `json:"range"`
	Severity int       `json:"severity"`
	Source   string    `json:"source"`
	Message  string    `json:"message"`
}

type diagRange struct {
	Start diagPosition `json:"start"`
	End   diagPosition `json:"end"`
}

type diagPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// publishDiagnosticsParams is the textDocument/publishDiagnostics notification body.
type publishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// computeDiagnostics scans all cached document lines for unknown req=/arch=/test-spec=
// values and returns a Diagnostic for each unknown ID.
func computeDiagnostics(uri string, lines []string, g *model.TraceabilityGraph) []Diagnostic {
	var diags []Diagnostic
	for lineNum, line := range lines {
		for _, m := range diagPattern.FindAllStringSubmatchIndex(line, -1) {
			attr := line[m[2]:m[3]]
			value := line[m[4]:m[5]]
			valStart := m[4]
			valEnd := m[5] // exclusive

			if isKnown(attr, value, g) {
				continue
			}

			diags = append(diags, Diagnostic{
				Range: diagRange{
					Start: diagPosition{Line: lineNum, Character: valStart},
					End:   diagPosition{Line: lineNum, Character: valEnd},
				},
				Severity: diagError,
				Source:   "req42-tracer",
				Message:  unknownMsg(attr, value),
			})
		}
	}
	return diags
}

// isKnown returns true if the given attr+value combination exists in the graph.
func isKnown(attr, value string, g *model.TraceabilityGraph) bool {
	switch attr {
	case "req":
		_, ok := g.Requirements[value]
		return ok
	case "arch":
		_, ok := g.ArchElements[value]
		return ok
	case "test-spec":
		_, ok := g.TestSpecs[value]
		return ok
	}
	return true // unknown attribute kind → don't flag
}

func unknownMsg(attr, value string) string {
	switch attr {
	case "req":
		return "Unknown requirement ID: " + value
	case "arch":
		return "Unknown architecture element: " + value
	case "test-spec":
		return "Unknown test spec ID: " + value
	}
	return "Unknown ID: " + value
}
