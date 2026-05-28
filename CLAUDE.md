# CLAUDE.md — req42-tracer

Requirements Tracing Tool für ASPICE PAM 4.0 (Go CLI).

## Kern-Konzept
- **Input**: AsciiDoc-Blöcke (`[req,id=...]`, `[arch,id=...]`, `[test-spec,...]`)
- **Architecture**: Bausteinsicht JSONC-Modell
- **Test-Results**: Dynamisch aus CI-Artifacts (JUnit XML, go-test JSON)
- **Output**: CLI (ASCII/Markdown) + HTML-Report (interaktiv)

## Projektstruktur
```
cmd/req42-tracer/       # CLI commands (trace, gaps, aspice, validate, watch, lsp)
internal/
  parser/              # AsciiDoc + Bausteinsicht-Parser
  model/               # Typen + Konfiguration
  graph/               # Traceability-Graph + Analyse
  aspice/              # ASPICE PAM 4.0 Prozess-Checker
  report/              # CLI + HTML Report-Generator
  testresult/          # Test-Result-Parser (JUnit, go-test)
  lsp/                 # LSP-Server (Minimal-MVP)
```

## Implementierungsplan
13 Schritte vom Skeleton bis zum LSP. Siehe `/home/coder/.claude/plans/kind-stargazing-torvalds.md`.

## Conventions
- Package: lowercase, kurz (z.B. `parser`, nicht `asciidoc_parser`)
- Typen: PascalCase (z.B. `Requirement`, `ArchElement`)
- Functions: camelCase (z.B. `parseAsciiDoc`, `buildGraph`)
- CLI: cobra-basiert mit subcommands
- Config: YAML (.req42.yaml)
- Tests: Dateiname + `_test.go` (z.B. `adoc.go` → `adoc_test.go`)

## Abhängigkeiten
- `github.com/spf13/cobra` — CLI Framework
- `github.com/fsnotify/fsnotify` — File-Watching
- `gopkg.in/yaml.v3` — YAML-Parsing

## Commands
- `req42-tracer trace` — Traceability Matrix
- `req42-tracer gaps` — Gap-Analyse
- `req42-tracer aspice` — ASPICE BP-Report
- `req42-tracer validate` — Model validation
- `req42-tracer watch --open` — Watch + live-reload
- `req42-tracer lsp` — LSP-Server

## Review Workflow
Siehe [`REVIEW.md`](REVIEW.md) für den vollständigen Code- und Security-Review-Prozess inkl. Inline-Kommentare im PR.
