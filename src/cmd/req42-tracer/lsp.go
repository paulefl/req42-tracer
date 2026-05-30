package main

import (
	"fmt"
	"os"

	"github.com/paulefl/req42-tracer/src/internal/lsp"
	"github.com/paulefl/req42-tracer/src/internal/model"
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

func runLspCmd(cmd *cobra.Command, _ []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	config, err := model.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[lsp] Warning: could not load config %s: %v\n", configPath, err)
	}
	if err := lsp.NewServer(config).Run(); err != nil {
		return fmt.Errorf("lsp server error: %w", err)
	}
	return nil
}
