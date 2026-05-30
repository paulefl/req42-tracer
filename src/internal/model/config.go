package model

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the .req42.yaml configuration file.
type Config struct {
	Projects       map[string]*ProjectConfig `yaml:"projects"`
	DefaultProject string                    `yaml:"default-project"` // optional; derived from first projects key if empty
	Bausteinsicht  struct {
		Model string `yaml:"model"`
	} `yaml:"bausteinsicht"`
	TestResults []struct {
		Format string `yaml:"format"` // junit, go-test-json
		Path   string `yaml:"path"`
	} `yaml:"test-results"`
	Rules      map[string]string `yaml:"rules"`       // error, warning, off
	RuleParams map[string]int    `yaml:"rule-params"` // numeric thresholds per rule
	ASPICE struct {
		AutoDerive    bool     `yaml:"auto-derive"`
		Processes     []string `yaml:"processes"`
		ProcessRules  map[string]map[string]string `yaml:"process-rules"`
	} `yaml:"aspice"`
	Reports struct {
		HTML struct {
			Output          string `yaml:"output"`
			IncludeGraph    bool   `yaml:"include-graph"`
			IncludeMatrix   bool   `yaml:"include-matrix"`
			IncludeASPICE   bool   `yaml:"include-aspice"`
			Theme           string `yaml:"theme"`
		} `yaml:"html"`
		CLI struct {
			Format string `yaml:"format"` // text, markdown, json
		} `yaml:"cli"`
	} `yaml:"reports"`
}

// ProjectConfig represents a single project in the configuration.
type ProjectConfig struct {
	Path string `yaml:"path"`
	Docs string `yaml:"docs"`
}

// LoadConfig loads the .req42.yaml configuration file.
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Set defaults
	if config.Projects == nil {
		config.Projects = make(map[string]*ProjectConfig)
	}
	if config.Rules == nil {
		config.Rules = make(map[string]string)
	}
	if config.ASPICE.ProcessRules == nil {
		config.ASPICE.ProcessRules = make(map[string]map[string]string)
	}

	// Default ASPICE processes if not specified
	if len(config.ASPICE.Processes) == 0 {
		config.ASPICE.Processes = []string{"SWE.1", "SWE.2", "SWE.3", "SWE.5"}
	}

	// Default report paths
	if config.Reports.HTML.Output == "" {
		config.Reports.HTML.Output = "reports/traceability-report.html"
	}
	if config.Reports.CLI.Format == "" {
		config.Reports.CLI.Format = "text"
	}

	return config, nil
}

// GetDefaultProject returns the configured default project name.
// Priority: explicit default-project field → first key in projects map → "software".
func (c *Config) GetDefaultProject() string {
	if c.DefaultProject != "" {
		return c.DefaultProject
	}
	for name := range c.Projects {
		return name
	}
	return "software"
}

// SetDefault sets a default value for a rule if not already set.
func (c *Config) SetDefault(rule, value string) {
	if _, exists := c.Rules[rule]; !exists {
		c.Rules[rule] = value
	}
}

// GetRule returns the rule level for a given rule name.
func (c *Config) GetRule(ruleName string) string {
	if level, exists := c.Rules[ruleName]; exists {
		return level
	}
	return "warning" // Default to warning if not specified
}
