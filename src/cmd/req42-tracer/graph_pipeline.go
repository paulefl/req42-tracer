package main

import (
	"fmt"
	"os"

	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
	"github.com/paulefl/req42-tracer/src/internal/parser"
)

// buildGraph parses requirements, architecture and the optional Bausteinsicht model,
// then derives ASPICE levels and builds trace links.
// reqDir and arcDir are the directories containing AsciiDoc sources.
// Returns the built traceability graph ready for analysis or reporting.
func buildGraph(config *model.Config, reqDir, arcDir, project string, verbose bool) (*model.TraceabilityGraph, error) {
	builder := graph.NewBuilder()

	if req, err := parser.ParseAllFromDir(reqDir, project); err == nil {
		if err := builder.MergeGraph(req); err != nil {
			return nil, fmt.Errorf("requirements merge from %s: %w", reqDir, err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed requirements from %s\n", reqDir)
		}
	} else if verbose {
		fmt.Fprintf(os.Stderr, "Warning: no requirements found in %s: %v\n", reqDir, err)
	}

	if arch, err := parser.ParseAllFromDir(arcDir, project); err == nil {
		if err := builder.MergeGraph(arch); err != nil {
			return nil, fmt.Errorf("architecture merge from %s: %w", arcDir, err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed architecture from %s\n", arcDir)
		}
	} else if verbose {
		fmt.Fprintf(os.Stderr, "Warning: no architecture found in %s: %v\n", arcDir, err)
	}

	if bPath := config.Bausteinsicht.Model; bPath != "" {
		loadBausteinsicht(builder, bPath, project, verbose)
	}

	builder.DeriveASPICELevels()
	if err := builder.BuildLinks(); err != nil {
		return nil, err
	}

	return builder.GetGraph(), nil
}
