# Plan: req42-tracer — MVP

## Context
Neues Go-CLI-Tool für Requirements Tracing nach ASPICE PAM 4.0.
AsciiDoc-Blöcke sind die primäre Source of Truth. Bausteinsicht liefert
Architecture-Elemente. CI-Artifacts liefern Test-Results dynamisch.

---

## Repository
- **Pfad**: `/home/coder/project/req42-tracer`
- **Go Module**: `github.com/paulefl/req42-tracer`
- **Go Version**: 1.24

---

## Kern-Konzept: AsciiDoc-Blöcke

Alle Artefakte werden als Block-Attribute in `.adoc` Dateien definiert:

```asciidoc
// REQ42 — Requirements
[req,id=SWE-001,version=2,priority=high,aspice=SWE.1,
 status=approved,reviewed-by=paulefl,reviewed-date=2026-05-07]
== Das System soll eine REST API bereitstellen

// ARC42 Kapitel 5, Level 1 → auto: aspice=SWE.2
[arch,id=comp.api,req=SWE-001,
 test-spec=spec.api.auth,
 impl=src/api/server.py]
== 5.1 API Component

// ARC42 Kapitel 5, Level 2 → auto: aspice=SWE.3 (hat parent)
[arch,id=comp.api.auth,parent=comp.api,req=SWE-002,
 impl=src/api/auth/handler.py]
=== 5.1.2 Authentication Handler

// Test-Spec
[test-spec,id=spec.api.auth,req=SWE-001,arch=comp.api]
== API Authentication Test Specification

// Test-Code (optional, primär via AsciiDoc; opt-in auch als Code-Annotation)
[test-code,id=test.api.auth,spec=spec.api.auth,
 file=tests/test_api_auth.py,fn=test_auth_endpoint]
```

---

## ASPICE-Mapping (automatisch aus ARC42-Hierarchie)

| Quelle | Regel | ASPICE-Prozess |
|--------|-------|----------------|
| ARC42 Kap. 3 | Kontextdiagramm | SYS.3 |
| ARC42 Kap. 5, kein `parent` | Whitebox Level 1 | SWE.2 |
| ARC42 Kap. 5, mit `parent` | Whitebox Level 2+ | SWE.3 |
| ARC42 Kap. 6 | Runtime View | SWE.2 |
| ARC42 Kap. 9 | ADRs | SWE.2 |
| Manuell überschreibbar | `aspice=SWE.1` im Block | — |

---

## Konfiguration (.req42.yaml)

```yaml
projects:
  system:
    path: ../system-repo/docs
    docs: docs/requirements/
  software:
    path: ./docs

bausteinsicht:
  model: architecture.jsonc

test-results:
  - format: junit
    path: reports/junit.xml
  - format: go-test-json
    path: reports/test-results.json

rules:
  undocumented-bausteinsicht-elements: warning  # error | warning | off
  missing-test-spec: error
  missing-review: warning
  stale-traces: warning
  missing-impl: warning
```

---

## Multi-Projekt & Cross-Repo

```asciidoc
[req,id=SWE-001,derives=system:SYS-001,version=1]
```

`system:SYS-001` → `system` ist Projekt-Key aus `.req42.yaml`.

---

## Versionierung & Stale-Detection

Wenn `version=N` hochgezählt wird, markiert das Tool alle abhängigen
Elemente die noch `req-version=N-1` referenzieren als stale:

```
⚠ STALE: comp.api.auth references SWE-001 v1 (current: v2)
  → Update req-version=2 in arc42/05-building-blocks.adoc:42
```

---

## Projektstruktur

```
req42-tracer/
├── cmd/req42-tracer/
│   ├── main.go
│   ├── root.go          # Cobra root, globale Flags (--model, --format, --config)
│   ├── trace.go         # cmd: trace — Traceability Matrix
│   ├── gaps.go          # cmd: gaps — Gap-Analyse
│   ├── aspice.go        # cmd: aspice — ASPICE BP-Coverage
│   ├── validate.go      # cmd: validate — Modell-Validierung
│   ├── watch.go         # cmd: watch — File-Watch Mode
│   └── lsp.go           # cmd: lsp — LSP-Server starten
├── internal/
│   ├── parser/
│   │   ├── adoc.go      # AsciiDoc Block-Scanner (regex-basiert)
│   │   └── bausteinsicht.go  # Bausteinsicht JSONC-Reader
│   ├── model/
│   │   ├── types.go     # Requirement, ArchElement, TestSpec, TestResult, TraceLink
│   │   └── config.go    # .req42.yaml laden
│   ├── graph/
│   │   ├── build.go     # Traceability-Graph aufbauen aus geparsten Blöcken
│   │   └── analysis.go  # Coverage, Orphans, Stale-Detection, Gap-Analyse
│   ├── aspice/
│   │   ├── processes.go # PAM 4.0 Prozess-Definitionen (SWE.1-SWE.6, SYS.1-SYS.3)
│   │   └── checker.go   # BP-Checks gegen Traceability-Graph
│   ├── report/
│   │   ├── table.go     # ASCII/Markdown Tabellenrenderer (CLI)
│   │   ├── html.go      # HTML-Report Generator
│   │   └── templates/   # HTML/CSS/JS Templates
│   │       └── report.html
│   ├── testresult/
│   │   ├── junit.go     # JUnit XML Parser
│   │   └── gotest.go    # go-test JSON Parser
│   └── lsp/
│       ├── server.go    # LSP-Server (stdio)
│       ├── complete.go  # Autocomplete für IDs
│       ├── hover.go     # Hover → Requirement-Text
│       └── diag.go      # Diagnostics (unbekannte IDs)
├── go.mod
├── Makefile
└── CLAUDE.md
```

---

## CLI Commands

```bash
# Traceability Matrix
req42-tracer trace --format=markdown

# Gap-Analyse
req42-tracer gaps
# → MISSING: SWE-002 hat kein test-spec
# → ORPHAN:  comp.legacy hat kein req

# ASPICE-Report
req42-tracer aspice
# → SWE.1 BP6 Testability: ⚠ 83% (10/12)
# → SWE.2 BP4 Traceability: ✓ 100%

# Validierung
req42-tracer validate

# Watch Mode
req42-tracer watch --open  # öffnet HTML-Report im Browser, live-reload

# LSP starten (für VS Code / IntelliJ)
req42-tracer lsp
```

---

## HTML-Report (interaktiv)

Drei Views, umschaltbar per Tab (wie unsere Test-Reports):

1. **Graph-View**: Klickbarer Dependency-Graph (D3.js oder Mermaid)
   - Nodes: req (blau), arch (grün), test-spec (orange), test-result (grau)
   - Edges: satisfies, implements, verifies, derives

2. **Matrix-View**: Traceability-Matrix
   - Rows = Requirements, Columns = Arch / Test-Specs
   - Farb-Coding: ✅ grün / ❌ rot / ⚠ orange (stale)

3. **ASPICE-Dashboard**: Pro Prozess eine Ampel + Coverage %
   - Konfigurierte Rules als Badges (error/warning/off)

Alle Views filterbar nach Projekt, ASPICE-Prozess, Status, Version.

---

## LSP (Minimal-MVP)

| Feature | Beschreibung |
|---------|-------------|
| Autocomplete | `req=` → listet alle bekannten Requirement-IDs |
| Hover | über `req=SWE-001` → zeigt Requirement-Text als Tooltip |
| Diagnostics | `req=UNKNOWN` → rote Unterstreichung + Message |
| Go-To-Definition | F12 auf `req=SWE-001` → springt zur Definitionszeile |

---

## ASPICE PAM 4.0 — MVP-Scope

| Prozess | Key BPs |
|---------|---------|
| SYS.2 | BP4 Traceability, BP5 Consistency |
| SYS.3 | BP3 Allocation |
| SWE.1 | BP5 Consistency, BP6 Testability, BP8 Bidirectional |
| SWE.2 | BP4 Traceability zu SWE.1 |
| SWE.3 | BP4 Traceability zu SWE.2 |
| SWE.5 | BP3 Traceability zu SWE.1 |

---

## Abhängigkeiten (go.mod)

```
github.com/spf13/cobra    v1.10.2  ← CLI
github.com/fsnotify/fsnotify v1.9.0 ← Watch-Mode
gopkg.in/yaml.v3          latest   ← .req42.yaml laden
```

LSP: eigene Implementierung über stdio (kein externes Framework für MVP).
HTML-Report: D3.js via CDN (kein Build-Step nötig).

---

## Implementierungsreihenfolge (MVP)

1. ✅ Repo-Skeleton (go.mod, Makefile, CLAUDE.md, root command)
2. `internal/model/types.go` — alle Typen
3. `internal/parser/adoc.go` — AsciiDoc Block-Scanner
4. `internal/parser/bausteinsicht.go` — JSONC-Reader
5. `internal/model/config.go` — .req42.yaml
6. `internal/graph/build.go` + `analysis.go`
7. `internal/aspice/processes.go` + `checker.go`
8. `internal/testresult/junit.go` + `gotest.go`
9. `internal/report/table.go` → CLI-Output (`trace`, `gaps`, `aspice`)
10. `internal/report/html.go` + Templates → HTML-Report
11. `watch.go` — File-Watch + live-reload
12. `internal/lsp/` — LSP Minimal-MVP
13. Beispiel-Projekt (REQ42 + ARC42 + Bausteinsicht Demo)

---

## Verifikation

```bash
make build

# Mit Demo-Projekt testen
cd examples/demo
req42-tracer validate
req42-tracer trace --format=markdown
req42-tracer gaps
req42-tracer aspice
req42-tracer watch --open   # Browser öffnet sich, live-reload

make test
```
