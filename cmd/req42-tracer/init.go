package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/paulefl/req42-tracer/internal/templates"
)

// NewInitCmd creates the `init` command for initializing new projects.
func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new REQ42 project",
		Long: `Initialize a new REQ42 + ARC42 project with templates.

This command creates the project skeleton including:
  - docs/requirements/req42.adoc (requirements document)
  - docs/arc42/arc42.adoc (architecture document)
  - architecture.jsonc (Bausteinsicht model)
  - .req42.yaml (configuration file)
  - .gitignore

Use --interactive=false with explicit flags for CI/CD automation.`,
		RunE: runInitCmd,
	}

	cmd.Flags().String("name", "", "Project name (default: req42-project)")
	cmd.Flags().String("module", "", "Go module path (default: github.com/user/project)")
	cmd.Flags().String("description", "", "Project description (default: REQ42 + ARC42 Project)")
	cmd.Flags().Bool("interactive", true, "Use interactive prompts (default: true)")

	return cmd
}

func runInitCmd(cmd *cobra.Command, args []string) error {
	interactive, _ := cmd.Flags().GetBool("interactive")

	var projectName, modulePath, description string
	var err error

	if interactive {
		projectName, modulePath, description, err = promptInteractive()
		if err != nil {
			return err
		}
	} else {
		projectName, _ = cmd.Flags().GetString("name")
		modulePath, _ = cmd.Flags().GetString("module")
		description, _ = cmd.Flags().GetString("description")

		if projectName == "" {
			projectName = "req42-project"
		}
		if modulePath == "" {
			modulePath = "github.com/user/project"
		}
		if description == "" {
			description = "REQ42 + ARC42 Project"
		}
	}

	return initializeProject(projectName, modulePath, description)
}

// promptInteractive prompts the user for project configuration interactively.
func promptInteractive() (name, module, description string, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Project name? [default: req42-project] ")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "req42-project"
	}

	fmt.Print("Module path? [default: github.com/user/project] ")
	module, _ = reader.ReadString('\n')
	module = strings.TrimSpace(module)
	if module == "" {
		module = "github.com/user/project"
	}

	fmt.Print("Description? [default: REQ42 + ARC42 Project] ")
	description, _ = reader.ReadString('\n')
	description = strings.TrimSpace(description)
	if description == "" {
		description = "REQ42 + ARC42 Project"
	}

	return name, module, description, nil
}

// initializeProject creates the project structure and processes templates.
func initializeProject(projectName, modulePath, description string) error {
	// Create directory structure
	dirs := []string{
		"docs/requirements",
		"docs/arc42",
		"reports",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Prepare substitution map
	placeholders := map[string]string{
		"{{PROJECT_NAME}}":       projectName,
		"{{PROJECT_NAME_SNAKE}}": toSnakeCase(projectName),
		"{{MODULE_PATH}}":        modulePath,
		"{{DESCRIPTION}}":        description,
	}

	// Process templates
	templateFiles := []struct {
		name     string
		dest     string
		isText   bool
	}{
		{"req42.adoc", "docs/requirements/req42.adoc", true},
		{"arc42.adoc", "docs/arc42/arc42.adoc", true},
		{"architecture.jsonc", "architecture.jsonc", true},
		{".req42.yaml", ".req42.yaml", true},
		{".gitignore", ".gitignore", true},
	}

	for _, tf := range templateFiles {
		content, err := templates.FS.ReadFile(tf.name)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", tf.name, err)
		}

		// Replace placeholders
		result := string(content)
		for placeholder, value := range placeholders {
			result = strings.ReplaceAll(result, placeholder, value)
		}

		// Write to destination
		if err := os.WriteFile(tf.dest, []byte(result), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", tf.dest, err)
		}

		fmt.Printf("✓ Created %s\n", tf.dest)
	}

	// Summary
	fmt.Println()
	fmt.Println("✨ Project initialized successfully!")
	fmt.Println()
	fmt.Printf("Next steps:\n")
	fmt.Printf("  1. Edit docs/requirements/req42.adoc to add your requirements\n")
	fmt.Printf("  2. Edit docs/arc42/arc42.adoc to document your architecture\n")
	fmt.Printf("  3. Update architecture.jsonc with your Bausteinsicht model\n")
	fmt.Printf("  4. Run: req42-tracer validate\n")
	fmt.Printf("  5. Run: req42-tracer trace\n")
	fmt.Printf("  6. Run: req42-tracer watch --open  (opens HTML report in browser)\n")

	return nil
}

// toSnakeCase converts a string to snake_case.
func toSnakeCase(s string) string {
	result := ""
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result += "_"
			}
			result += string(r + 32) // Convert to lowercase
		} else if r == ' ' {
			result += "_"
		} else {
			result += string(r)
		}
	}
	return strings.ToLower(result)
}
