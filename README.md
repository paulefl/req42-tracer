# req42-tracer — Requirements Tracing Tool

A Go CLI tool for tracing requirements across AsciiDoc documentation, architecture models (Bausteinsicht), and test specifications with ASPICE PAM 4.0 compliance validation.

## MVP Features

### Core Commands

- **`req42-tracer init`** — Initialize new project with templates (interactive or automated)
- **`req42-tracer trace`** — Display traceability matrix showing requirement → architecture → test links
- **`req42-tracer gaps`** — Analyze and report orphan artifacts and missing implementations
- **`req42-tracer aspice`** — Validate ASPICE PAM 4.0 compliance, show process coverage
- **`req42-tracer validate`** — Check project structure and validate all references

### Architecture

- **AsciiDoc Parser** — Extracts [req], [arch], [test-spec] blocks with attributes
- **Bausteinsicht Integration** — Loads architecture models as JSONC with nested hierarchy
- **Traceability Graph** — Bidirectional dependency graph with explicit and name-based linking
- **ASPICE Checker** — Validates SWE.1-6, SYS.1-3 processes with best practice checks
- **Test Result Parsers** — JUnit XML and go-test JSON support with platform detection
- **Report Generator** — Text, Markdown, and JSON output formats

### Hybrid Test-Result Tracing

Combines three approaches for matching test results to specifications:
1. **Explicit annotations** in test code comments (primary)
2. **Name-based matching** using heuristics (fallback)
3. **Custom metadata** in test reports (optional)

## Quick Start

```bash
# Initialize a project
req42-tracer init --name=MyProject --module=github.com/user/myproject --interactive=false

# Validate structure
req42-tracer validate

# Show traceability
req42-tracer trace --format=markdown

# Analyze gaps
req42-tracer gaps

# Check ASPICE compliance
req42-tracer aspice
```

## Project Structure

```
internal/
  ├── model/       # Type definitions, config loader
  ├── parser/      # AsciiDoc + Bausteinsicht parsers
  ├── graph/       # Traceability builder and analyzer
  ├── aspice/      # ASPICE PAM 4.0 registry and checker
  ├── testresult/  # JUnit/go-test JSON loaders
  ├── report/      # CLI report generators
  └── templates/   # Project initialization templates

cmd/req42-tracer/
  ├── root.go      # CLI root command
  ├── init.go      # Project initialization
  ├── trace.go     # Traceability matrix
  ├── gaps.go      # Gap analysis
  ├── aspice.go    # ASPICE compliance
  └── validate.go  # Validation
```

## Roadmap (Phase 2)

- [ ] HTML Report Generator with interactive Graph/Matrix/ASPICE views
- [ ] Watch Mode with live-reload browser preview
- [ ] LSP Server for IDE integration (autocomplete, hover, diagnostics)
- [ ] Advanced filtering and search
- [ ] Custom validation rules

## Implementation Status

✅ Complete:
- Core type system
- AsciiDoc and Bausteinsicht parsing
- Traceability graph construction
- ASPICE PAM 4.0 validation
- Test result parsing
- CLI reports (text/markdown/json)
- Project initialization

⏳ Pending:
- HTML report generation
- Watch mode
- LSP server
