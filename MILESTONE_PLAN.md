# req42-tracer — Milestone Plan

**GitHub Issues & Milestones**: https://github.com/paulefl/req42-tracer

---

## MVP ✅ (Complete)

Core CLI functionality for requirements tracing.

**Completed Features**:
- ✅ `req42-tracer init` — Project initialization with templates
- ✅ `req42-tracer trace` — Traceability matrix (text/markdown/json)
- ✅ `req42-tracer gaps` — Gap analysis report
- ✅ `req42-tracer aspice` — ASPICE PAM 4.0 compliance check
- ✅ `req42-tracer validate` — Reference validation
- ✅ AsciiDoc parser ([req], [arch], [test-spec] blocks)
- ✅ Bausteinsicht JSONC model loader
- ✅ Traceability graph construction
- ✅ Test-result parsers (JUnit, go-test JSON)

**Story Points**: ~50 SP

---

## Phase 2: HTML & Watch 🔄 (In Progress)

Interactive HTML reports with live-reload.

### Issues Created

| # | Title | Story Points | Status |
|---|-------|------|--------|
| [#1](https://github.com/paulefl/req42-tracer/issues/1) | HTML Report with D3.js Graph | 8 | 🟡 Ready |
| [#2](https://github.com/paulefl/req42-tracer/issues/2) | Matrix View Tab | 5 | 🟡 Ready |
| [#3](https://github.com/paulefl/req42-tracer/issues/3) | ASPICE Dashboard Tab | 5 | 🟡 Ready |
| [#4](https://github.com/paulefl/req42-tracer/issues/4) | Watch Mode with Live-Reload | 5 | 🟡 Ready |

**Total**: 23 story points

### Acceptance Criteria

- [ ] `req42-tracer` generates `reports/traceability-report.html`
- [ ] Three interactive tabs: Graph, Matrix, ASPICE
- [ ] D3.js dependency graph (clickable nodes, filterable edges)
- [ ] `req42-tracer watch --open` opens browser and auto-reloads on file changes
- [ ] Demo project shows complete traceability visualization

---

## Phase 3: IDE Integration 🎯

Language Server Protocol for VS Code / IntelliJ.

### Issues Created

| # | Title | Story Points | Status |
|---|-------|------|--------|
| [#5](https://github.com/paulefl/req42-tracer/issues/5) | LSP Server Skeleton | 3 | 🟡 Ready |
| [#6](https://github.com/paulefl/req42-tracer/issues/6) | Autocomplete req/arch/test IDs | 3 | 🟡 Ready |
| [#7](https://github.com/paulefl/req42-tracer/issues/7) | Hover Tooltips | 2 | 🟡 Ready |
| [#8](https://github.com/paulefl/req42-tracer/issues/8) | Diagnostics & Error Underlines | 3 | 🟡 Ready |
| [#9](https://github.com/paulefl/req42-tracer/issues/9) | Go-to-Definition (F12) | 2 | 🟡 Ready |

**Total**: 13 story points

### Acceptance Criteria

- [ ] `req42-tracer lsp` starts server on stdio
- [ ] VS Code extension can connect
- [ ] Full autocomplete + hover + diagnostics + go-to-def workflow
- [ ] Works in .adoc and .jsonc files

---

## Phase 4: Dogfooding 🐕

Apply the tool to itself — define req42-tracer's own requirements & architecture.

### Issues Created

| # | Title | Story Points | Status |
|---|-------|------|--------|
| [#10](https://github.com/paulefl/req42-tracer/issues/10) | Document req42-tracer Requirements | 3 | 🟡 Ready |
| [#11](https://github.com/paulefl/req42-tracer/issues/11) | Create ADRs | 3 | 🟡 Ready |
| [#12](https://github.com/paulefl/req42-tracer/issues/12) | Complete arc42.adoc with [arch] blocks | 4 | 🟡 Ready |
| [#13](https://github.com/paulefl/req42-tracer/issues/13) | Unit Tests (80%+ coverage) | 8 | 🟡 Ready |

**Total**: 18 story points

### Acceptance Criteria

- [ ] docs/requirements/req42.adoc has functional & non-functional requirements
- [ ] docs/arc42/ADRs/ contains 4+ ADRs with decisions
- [ ] docs/arc42/arc42.adoc fully documented with [arch] blocks
- [ ] `req42-tracer trace` shows complete traceability for req42-tracer itself
- [ ] Test coverage: parser ≥80%, graph ≥80%, aspice ≥75%
- [ ] `req42-tracer aspice` shows 100% (or close) compliance for Phase 4 requirements

---

## Phase 5: Production Ready 🚀

Performance, documentation, CI/CD, advanced features.

### Issues Created

| # | Title | Story Points | Status |
|---|-------|------|--------|
| [#14](https://github.com/paulefl/req42-tracer/issues/14) | Performance Optimization | 5 | 🟡 Ready |
| [#15](https://github.com/paulefl/req42-tracer/issues/15) | User Documentation | 5 | 🟡 Ready |
| [#16](https://github.com/paulefl/req42-tracer/issues/16) | CI/CD Pipeline (GitHub Actions) | 4 | 🟡 Ready |
| [#17](https://github.com/paulefl/req42-tracer/issues/17) | Custom Validation Rules | 5 | 🟡 Ready |

**Total**: 19 story points

### Acceptance Criteria

- [ ] Trace command: <2s for 1000+ requirements
- [ ] Comprehensive user guide with examples
- [ ] GitHub Actions: build, test, lint, release binaries
- [ ] Custom rules in .req42.yaml with configurable severity

---

## Timeline Estimate

| Phase | Story Points | Estimated Duration | Status |
|-------|-----|---------|--------|
| MVP | 50 | ✅ Complete | Done |
| Phase 2 | 23 | ~2-3 weeks | Ready to start |
| Phase 3 | 13 | ~2 weeks | Ready after Phase 2 |
| Phase 4 | 18 | ~2 weeks (parallel) | Can start anytime |
| Phase 5 | 19 | ~2-3 weeks | Ready after Phase 2 |

**Total**: 123 story points (~12-16 weeks estimated)

---

## How to Start Contributing

1. Pick an issue from the milestone
2. Assign to yourself
3. Create a feature branch: `git checkout -b <issue-#>-short-name`
4. Follow commit conventions: `feat(module): description`
5. Create a PR linked to the issue
6. Ensure: tests pass, coverage maintained, linter clean

---

## Legend

- 🟡 Ready — Issue defined, ready to pick up
- 🟢 In Progress — Actively being worked on
- ✅ Done — Completed and merged

---

**GitHub Repo**: https://github.com/paulefl/req42-tracer  
**Last Updated**: 2026-05-07
