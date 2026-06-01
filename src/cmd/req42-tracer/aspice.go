package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/src/internal/aspice"
	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
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

	config, err := model.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	if verbose {
		fmt.Fprintf(os.Stderr, "Loaded config from %s\n", configPath)
	}

	g, err := buildGraph(config,
		"project/req42-tracer/docs/requirements",
		"project/req42-tracer/docs/arc42",
		config.GetDefaultProject(), verbose)
	if err != nil {
		return err
	}

	analyzer := graph.NewAnalyzer(g)
	checker := aspice.NewChecker(analyzer, config)
	rpt := checker.CheckCompliance()

	fmt.Println("ASPICE PAM 4.0 Compliance Report")
	fmt.Printf("Overall Coverage: %.1f%%\n\n", rpt.Overall)

	for processID, results := range rpt.Processes {
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
