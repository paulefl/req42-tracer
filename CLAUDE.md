# CLAUDE.md — req42-tracer

Requirements Tracing Tool für ASPICE PAM 4.0 (Go CLI).

## Kern-Konzept
- **Input**: AsciiDoc-Blöcke (`[req,id=...]`, `[arch,id=...]`, `[test-spec,...]`)
- **Architecture**: Bausteinsicht JSONC-Modell
- **Test-Results**: Dynamisch aus CI-Artifacts (JUnit XML, go-test JSON)
- **Output**: CLI (ASCII/Markdown) + HTML-Report (interaktiv)

## Projektstruktur
```
src/
  cmd/req42-tracer/    # CLI commands (trace, gaps, aspice, validate, watch, lsp)
  internal/
    parser/            # AsciiDoc + Bausteinsicht-Parser
    model/             # Typen + Konfiguration
    graph/             # Traceability-Graph + Analyse
    aspice/            # ASPICE PAM 4.0 Prozess-Checker
    report/            # CLI + HTML Report-Generator
    testresult/        # Test-Result-Parser (JUnit, go-test)
    lsp/               # LSP-Server (Minimal-MVP)
    validation/        # Custom Validation Rules Engine
project/req42-tracer/
  docs/                # arc42, requirements, user-guide
  scripts/             # generate-diagrams.sh etc.
  tools/bausteinsicht/ # Bausteinsicht CLI v1.1.0
```

## Bausteinsicht
- **Tool:** `project/req42-tracer/tools/bausteinsicht/bausteinsicht` (Linux amd64)
- **Schema:** v1.1.0 — `architecture.jsonc` enthält `$schema`-Referenz
- **Validieren:** `./project/req42-tracer/tools/bausteinsicht/bausteinsicht validate --model project/req42-tracer/architecture.jsonc`

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

## Setup
Siehe [`SETUP.md`](project/req42-tracer/SETUP.md) für die einmalige Einrichtung der Entwicklungsumgebung (Repos klonen, Claude Skill Library installieren, gh CLI einrichten).

## Rollen
Siehe [`ROLES.md`](project/req42-tracer/ROLES.md) für die vollständige Rollendefinition.

- **Implementierung:** `dev-paul-fleischmann` — Feature-Branches, Commits, PRs öffnen
- **Review & Merge:** `paulefl` — Code/Security-Review, Approve, Merge

### Pflicht-Checkliste vor jedem Schritt

**Vor Implementierung (Commit/Push/PR erstellen):**
```bash
gh auth switch --user dev-paul-fleischmann
git config user.name "dev-paul-fleischmann"
git config user.email "dev@paul-fleischmann.com"
```

**Vor Review/Merge:**
```bash
gh auth switch --user paulefl
```

Ein pre-commit Hook in `.git/hooks/pre-commit` blockiert Commits mit falscher Identität.
`git config user.name` muss `dev-paul-fleischmann` sein — sonst bricht der Commit ab.

Ein `commit-msg` Hook in `.git/hooks/commit-msg` erzwingt Issue-Referenz im Format:
```
Closes #<NNN>: <beschreibung>
Fixes #<NNN>: <beschreibung>
Refs #<NNN>: <beschreibung>
```
Merge-Commits, fixup! und squash! sind ausgenommen.

PR erstellen: `gh pr create --assignee dev-paul-fleischmann --reviewer paulefl`

## Commit-Konventionen

Jeder Commit der eine Issue-Arbeit repräsentiert **muss** eine Issue-Referenz enthalten:

```
Closes #91: ci add deadcode job
Fixes #75: wire testresult.LoadAll into trace command
Refs #93: partial progress on BausteinsichtParser
```

| Keyword | Effekt beim PR-Merge |
|---|---|
| `Closes #N` / `Fixes #N` / `Resolves #N` | Issue wird automatisch geschlossen |
| `Refs #N` | Issue bleibt offen, Commit erscheint als Referenz |

Merge-Commits, fixup! und squash! brauchen keine Referenz.

## Test Conventions
Siehe [`TESTS.md`](project/req42-tracer/TESTS.md) für die vollständigen Test-Konventionen.

Kurzfassung: Jeder Test braucht einen `[test-spec]`-Block mit `req=` und `aspice=` Linkage direkt oberhalb der Testfunktion.

### Coverage-Ziele (Phase 4)
| Paket | Ziel |
|---|---|
| `src/internal/parser` | ≥ 80% |
| `src/internal/graph` | ≥ 80% |
| `src/internal/aspice` | ≥ 75% |
| `src/internal/report` | ≥ 70% |
| `src/internal/model` | ≥ 60% |
| `src/internal/testresult` | ≥ 70% |

Coverage prüfen: `go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out`

## Review Workflow
Siehe [`REVIEW.md`](project/req42-tracer/REVIEW.md) für den vollständigen Code- und Security-Review-Prozess inkl. Inline-Kommentare im PR.
