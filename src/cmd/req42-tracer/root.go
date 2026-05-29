package main

import (
	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags "-X main.version=<tag>"
var version = "dev"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "req42-tracer",
		Short: "Requirements Tracing Tool for ASPICE PAM 4.0",
		Long: `req42-tracer is a CLI tool for tracing requirements across AsciiDoc
documentation, architecture models (Bausteinsicht), and test specifications.
It supports ASPICE PAM 4.0 process validation and generates interactive reports.`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Global flags
	cmd.PersistentFlags().String("config", ".req42.yaml", "Configuration file path")
	cmd.PersistentFlags().String("format", "text", "Output format: text, markdown, json, html")
	cmd.PersistentFlags().Bool("verbose", false, "Verbose output")

	// Subcommands
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newTraceCmd())
	cmd.AddCommand(newGapsCmd())
	cmd.AddCommand(newAspiceCmd())
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newWatchCmd())
	cmd.AddCommand(newLspCmd())

	return cmd
}
