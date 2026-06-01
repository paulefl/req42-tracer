package main

import (
	"fmt"
	"os"

	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
	"github.com/paulefl/req42-tracer/src/internal/parser"
	"github.com/paulefl/req42-tracer/src/internal/testresult"
)

// buildGraph parses requirements, architecture, Bausteinsicht model, and Go test
// annotations, then derives ASPICE levels and builds trace links.
// reqDir and arcDir are the directories containing AsciiDoc sources.
// goSrcDir is the root directory scanned for *_test.go files (empty = skip).
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

	// Parse Go test files for [test-spec] annotations → explicit TestCode entries
	if goSrc := config.GoSrcDir; goSrc != "" {
		if goGraph, err := parser.ParseGoTestFiles(goSrc, project); err == nil {
			if mergeErr := builder.MergeGraph(goGraph); mergeErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: Go test code merge: %v\n", mergeErr)
			} else if verbose {
				fmt.Fprintf(os.Stderr, "Parsed Go test annotations from %s\n", goSrc)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Warning: could not parse Go test files: %v\n", err)
		}
	}

	// Load test results before BuildLinks so name-based linking can match them.
	if err := testresult.LoadAll(builder.GetGraph(), config); err != nil && verbose {
		fmt.Fprintf(os.Stderr, "Warning: could not load test results: %v\n", err)
	}

	builder.DeriveASPICELevels()
	if err := builder.BuildLinks(); err != nil {
		return nil, err
	}

	return builder.GetGraph(), nil
}
