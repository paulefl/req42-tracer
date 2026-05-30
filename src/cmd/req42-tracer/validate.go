package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
	"github.com/paulefl/req42-tracer/src/internal/parser"
	"github.com/paulefl/req42-tracer/src/internal/validation"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate project structure and custom rules",
		Long: `Validate that the project is correctly structured, all references
are valid, and custom validation rules from .req42.yaml pass.`,
		RunE: runValidateCmd,
	}
	return cmd
}

func runValidateCmd(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	verbose, _ := cmd.Flags().GetBool("verbose")

	cfg, err := model.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if verbose {
		fmt.Fprintf(os.Stderr, "Loaded config from %s\n", configPath)
	}

	builder := graph.NewBuilder()
	project := cfg.GetDefaultProject()

	if req, err := parser.ParseAllFromDir("project/req42-tracer/docs/requirements", project); err == nil {
		if err := builder.MergeGraph(req); err != nil {
			return err
		}
	} else if verbose {
		fmt.Fprintf(os.Stderr, "No requirements found: %v\n", err)
	}

	if arch, err := parser.ParseAllFromDir("project/req42-tracer/docs/arc42", project); err == nil {
		if err := builder.MergeGraph(arch); err != nil {
			return err
		}
	} else if verbose {
		fmt.Fprintf(os.Stderr, "No architecture found: %v\n", err)
	}

	// Load Bausteinsicht model if configured
	if bPath := cfg.Bausteinsicht.Model; bPath != "" {
		loadBausteinsicht(builder, bPath, project, verbose)
	}

	builder.DeriveASPICELevels()
	if err := builder.BuildLinks(); err != nil {
		return err
	}

	g := builder.GetGraph()
	analyzer := graph.NewAnalyzer(g)

	// Built-in reference validation
	refErrors := analyzer.ValidateReferences()

	// Custom rules from config
	engine := validation.NewRuleEngine(cfg, analyzer)
	ruleResults := engine.Run()
	numErrors, numWarnings := validation.TotalViolations(ruleResults)

	// --- Output ---
	if len(refErrors) > 0 {
		fmt.Println("❌ Reference errors:")
		for _, e := range refErrors {
			fmt.Printf("  ❌ [ERROR] %s\n", e)
		}
	}

	if output := validation.FormatResults(ruleResults); output != "" {
		fmt.Print(output)
	}

	totalErrors := len(refErrors) + numErrors
	if totalErrors == 0 && numWarnings == 0 {
		fmt.Println("✅ Validation successful!")
	} else if totalErrors == 0 {
		fmt.Printf("✅ Validation passed with %d warning(s)\n", numWarnings)
	} else {
		fmt.Printf("❌ Validation failed: %d error(s), %d warning(s)\n", totalErrors, numWarnings)
	}

	fmt.Printf("  Requirements: %d | Architecture: %d | Test Specs: %d | Links: %d\n",
		len(g.Requirements), len(g.ArchElements), len(g.TestSpecs), len(g.Links))

	if totalErrors > 0 {
		return fmt.Errorf("validation failed with %d error(s)", totalErrors)
	}
	return nil
}
