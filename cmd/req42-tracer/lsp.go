package main

import (
	"fmt"

	"github.com/paulefl/req42-tracer/internal/lsp"
	"github.com/spf13/cobra"
)

func newLspCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lsp",
		Short: "Start LSP server (stdio)",
		Long: `Start a Language Server Protocol server on stdio.

Connect your editor to req42-tracer lsp to get autocompletion,
hover tooltips, diagnostics, and go-to-definition for req42 blocks
in .adoc files.`,
		RunE: runLspCmd,
	}
}

func runLspCmd(_ *cobra.Command, _ []string) error {
	if err := lsp.NewServer().Run(); err != nil {
		return fmt.Errorf("lsp server error: %w", err)
	}
	return nil
}
