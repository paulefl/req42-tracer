package main

import (
	"fmt"
	"os"

	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/parser"
)

// loadBausteinsicht parses the Bausteinsicht JSONC model and merges it into builder.
// Parse errors are always reported (not gated on --verbose).
// JSONC elements whose IDs conflict with already-loaded arch elements are skipped
// with a warning instead of aborting the command.
func loadBausteinsicht(builder *graph.Builder, bPath, project string, verbose bool) {
	bParser := parser.NewBausteinsichtParser(bPath)
	bGraph, err := bParser.Parse(project)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load Bausteinsicht model: %v\n", err)
		return
	}

	// Pre-filter elements that would conflict with already-loaded arch elements
	// to avoid a hard MergeGraph failure. JSONC is lower priority than AsciiDoc.
	current := builder.GetGraph()
	skipped := 0
	for id := range bGraph.ArchElements {
		if _, exists := current.ArchElements[id]; exists {
			delete(bGraph.ArchElements, id)
			skipped++
		}
	}
	if skipped > 0 {
		fmt.Fprintf(os.Stderr, "Warning: skipped %d Bausteinsicht element(s) with duplicate IDs (AsciiDoc takes precedence)\n", skipped)
	}

	if err := builder.MergeGraph(bGraph); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Bausteinsicht merge error: %v\n", err)
		return
	}
	if verbose {
		fmt.Fprintf(os.Stderr, "Parsed Bausteinsicht model from %s\n", bPath)
	}
}
