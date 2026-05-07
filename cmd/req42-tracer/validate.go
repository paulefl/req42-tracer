package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/internal/model"
	"github.com/paulefl/req42-tracer/internal/parser"
	"github.com/paulefl/req42-tracer/internal/graph"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate project structure and references",
		Long: `Validate that the project is correctly structured and all
references are valid (requirements exist, architecture IDs match, etc.).`,
		RunE: runValidateCmd,
	}

	return cmd
}

func runValidateCmd(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Load configuration (for future use with advanced options)
	_, err := model.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Loaded config from %s\n", configPath)
	}

	// Build traceability graph
	builder := graph.NewBuilder()

	// Parse requirements
	if req, err := parser.ParseAllFromDir("docs/requirements", "software"); err == nil {
		if err := builder.MergeGraph(req); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed requirements\n")
		}
	} else if verbose {
		fmt.Fprintf(os.Stderr, "No requirements found: %v\n", err)
	}

	// Parse architecture
	if arch, err := parser.ParseAllFromDir("docs/arc42", "software"); err == nil {
		if err := builder.MergeGraph(arch); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed architecture\n")
		}
	} else if verbose {
		fmt.Fprintf(os.Stderr, "No architecture found: %v\n", err)
	}

	// Derive ASPICE levels
	builder.DeriveASPICELevels()

	// Build trace links
	if err := builder.BuildLinks(); err != nil {
		return err
	}

	// Get final graph
	g := builder.GetGraph()

	// Validate references
	analyzer := graph.NewAnalyzer(g)
	errors := analyzer.ValidateReferences()

	// Display results
	if len(errors) == 0 {
		fmt.Println("✅ Project validation successful!")
		fmt.Printf("  Requirements: %d\n", len(g.Requirements))
		fmt.Printf("  Architecture Elements: %d\n", len(g.ArchElements))
		fmt.Printf("  Test Specifications: %d\n", len(g.TestSpecs))
		fmt.Printf("  Trace Links: %d\n", len(g.Links))
		return nil
	}

	fmt.Println("❌ Validation errors found:")
	for _, err := range errors {
		fmt.Printf("  - %s\n", err)
	}

	return fmt.Errorf("validation failed with %d errors", len(errors))
}
