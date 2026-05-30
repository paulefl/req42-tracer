package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/src/internal/aspice"
	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
	"github.com/paulefl/req42-tracer/src/internal/parser"
	"github.com/paulefl/req42-tracer/src/internal/testresult"
)

func newAspiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aspice",
		Short: "Validate ASPICE PAM 4.0 compliance",
		Long: `Validate compliance against ASPICE PAM 4.0 processes and best practices.
Displays coverage percentages and identifies gaps in each process.`,
		RunE: runAspiceCmd,
	}

	return cmd
}

func runAspiceCmd(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Load configuration
	config, err := model.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Loaded config from %s\n", configPath)
	}

	// Build traceability graph
	builder := graph.NewBuilder()

	// Parse requirements
	if req, err := parser.ParseAllFromDir("project/req42-tracer/docs/requirements", "software"); err == nil {
		if err := builder.MergeGraph(req); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed requirements\n")
		}
	}

	// Parse architecture
	if arch, err := parser.ParseAllFromDir("project/req42-tracer/docs/arc42", "software"); err == nil {
		if err := builder.MergeGraph(arch); err != nil {
			return err
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed architecture\n")
		}
	}

	// Derive ASPICE levels
	builder.DeriveASPICELevels()

	// Build trace links
	if err := builder.BuildLinks(); err != nil {
		return err
	}

	// Get final graph
	g := builder.GetGraph()

	// Load test results from CI artifacts
	if err := testresult.LoadAll(g, config); err != nil && verbose {
		fmt.Fprintf(os.Stderr, "Warning: could not load test results: %v\n", err)
	}

	// Validate ASPICE compliance
	analyzer := graph.NewAnalyzer(g)
	checker := aspice.NewChecker(analyzer, config)
	report := checker.CheckCompliance()

	// Display report
	fmt.Println("ASPICE PAM 4.0 Compliance Report")
	fmt.Printf("Overall Coverage: %.1f%%\n\n", report.Overall)

	for processID, results := range report.Processes {
		totalCoverage := 0.0
		for _, result := range results {
			totalCoverage += result.Coverage
		}
		avgCoverage := totalCoverage / float64(len(results))

		status := "✅"
		if avgCoverage < 100 {
			status = "⚠"
		}
		if avgCoverage < 50 {
			status = "❌"
		}

		fmt.Printf("%s %s: %.1f%%\n", status, processID, avgCoverage)

		for _, result := range results {
			if result.Coverage < 100 {
				fmt.Printf("  - %s: %.1f%%\n", result.BP.Title, result.Coverage)
				for _, gap := range result.Gaps {
					fmt.Printf("    → %s\n", gap)
				}
			}
		}
		fmt.Println()
	}

	return nil
}
