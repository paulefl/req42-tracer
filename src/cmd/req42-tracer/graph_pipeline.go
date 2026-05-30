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
			return nil, err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed requirements from %s\n", reqDir)
		}
	}

	if arch, err := parser.ParseAllFromDir(arcDir, project); err == nil {
		if err := builder.MergeGraph(arch); err != nil {
			return nil, err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed architecture from %s\n", arcDir)
		}
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
