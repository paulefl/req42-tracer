# req42-tracer — Ideas & Future Directions

Brainstorming-Sammlung für Features, Verbesserungen und Integrationen.
Nicht priorisiert — dient als Diskussionsbasis.

---

## 🔍 Parser & Input-Quellen

| # | Idee | Beschreibung |
|---|------|-------------|
| P-01 | **Rust-Support** | `[req,id=...]` in `.rs`-Kommentaren parsen |
| P-02 | **Python-Support** | Docstrings mit `[req,id=...]` erkennen |
| P-03 | **C/C++-Support** | `[req,id=...]` in `.c/.h/.cpp`-Kommentaren (siehe Issue #106) |
| P-04 | **Doxygen-Integration** | `@req`/`@arch`-Tags aus Doxygen XML extrahieren (Issue #111) |
| P-05 | **Swagger/OpenAPI-Parser** | API-Endpoints als Arch-Elemente importieren |
| P-06 | **Confluence-Import** | Requirements direkt aus Confluence-Pages lesen |
| P-07 | **JIRA-Import** | Issues als Requirements importieren |
| P-08 | **ReqIF/DOORS-Import** | Standard-Austauschformat einlesen (Issue #78) |
| P-09 | **Excel/CSV-Import** | Requirements aus Tabellen importieren |
| P-10 | **Postman/Newman** | API-Tests als Test-Results einlesen |
| P-11 | **Playwright/Cypress** | E2E-Test-Ergebnisse tracen |
| P-12 | **Robot Framework** | Test-Reports parsen |
| P-13 | **Git-Commits als Trace-Events** | `Closes #REQ-001` in Commits → automatische Traceability |

---

## 📊 Report & Visualisierung

| # | Idee | Beschreibung |
|---|------|-------------|
| R-01 | **Coverage-Dashboard** | Kombinierter Test-Report analog Bausteinsicht `test-report-combined` (Issue #169) |
| R-02 | **Sankey-Diagramm** | Traceability-Flow von Req → Arch → Test → Result visualisieren |
| R-03 | **Heatmap** | Welche Komponenten haben die wenigste Coverage (farbkodiert) |
| R-04 | **Sunburst-Chart** | Hierarchische Bausteinsicht mit Coverage-Overlay |
| R-05 | **Timeline-View** | Wann wurde welche Anforderung umgesetzt/getestet |
| R-06 | **Dependency-Graph** | Welche Anforderungen hängen voneinander ab |
| R-07 | **Diff-View** | Traceability-Matrix zwischen zwei Git-Commits vergleichen (Issue #83) |
| R-08 | **Coverage-Badges** | SVG-Badges für README (z.B. `Requirements: 100% ✓`) (Issue #86) |
| R-09 | **PDF-Export** | ASPICE-Audit-Report als PDF (Issue #85) |
| R-10 | **Word/DOCX-Export** | Traceability-Matrix als Word-Dokument |
| R-11 | **LaTeX-Export** | Für wissenschaftliche/formale Dokumentation |
| R-12 | **Dark Mode** | Für den HTML-Report |
| R-13 | **Druckansicht** | Optimiertes CSS für `@media print` |
| R-14 | **Shareable Links** | Report-State (Filter, Tab) in URL kodieren (Deep-Links) |
| R-15 | **Embedded Mermaid** | Automatisch Traceability-Diagramme in `.adoc` einbetten |
| R-16 | **Trend-Dashboard** | Coverage-Verlauf über mehrere CI-Runs hinweg (Issue #84) |

---

## 🤖 KI / Automatisierung

| # | Idee | Beschreibung |
|---|------|-------------|
| A-01 | **AI Test-Spec-Vorschläge** | Aus einer Anforderung automatisch `[test-spec]`-Blöcke generieren (Issue #96) |
| A-02 | **AI Gap-Erklärung** | KI erklärt warum eine Lücke existiert und schlägt Fix vor |
| A-03 | **AI Requirements-Review** | Prüft Anforderungen auf Vollständigkeit und Testbarkeit |
| A-04 | **Duplicate-Detection** | Semantisch ähnliche Anforderungen erkennen (Issue #80) |
| A-05 | **AI-Projektbericht** | "Dein Projekt hat 87% Coverage. Die kritischsten Lücken sind..." |
| A-06 | **Natural Language Query** | "Welche Anforderungen sind nicht getestet?" als Freitext stellen |
| A-07 | **Auto-Link-Vorschläge** | KI schlägt fehlende `req=`/`arch=` Links vor |
| A-08 | **Change-Impact-Analyse** | Bei Änderung an REQ-001: welche Tests müssen angepasst werden (Issue #94) |

---

## 🔗 Integrationen & Ecosystem

| # | Idee | Beschreibung |
|---|------|-------------|
| I-01 | **GitHub Actions Output** | Traceability-Matrix als PR-Kommentar posten |
| I-02 | **GitHub Issues ↔ Requirements** | Bidirektionale Synchronisation (Issue #77) |
| I-03 | **JIRA-Ticket-Link** | `[req,jira=PROJ-123]` im HTML klickbar machen (Issue #76) |
| I-04 | **Linear-Integration** | Requirements mit Linear-Issues verlinken |
| I-05 | **Notion-Export** | Traceability-Matrix in Notion-Database schreiben |
| I-06 | **Slack-Notifications** | Alert wenn Coverage unter Schwellenwert fällt |
| I-07 | **VS Code Extension** | Tree-View der Traceability im Editor (Issue #87) |
| I-08 | **IntelliJ Plugin** | Analog VS Code Extension |
| I-09 | **Neovim Plugin** | LSP ist vorhanden, Plugin-Wrapper drauf |
| I-10 | **Pre-commit Hook** | Automatisch Traceability prüfen vor jedem Commit |
| I-11 | **GitHub Bot** | PR-Review automatisch auf Traceability-Vollständigkeit prüfen |
| I-12 | **Confluence-Export** | Traceability-Matrix nach Confluence publizieren (Issue #79) |

---

## 🏗️ Architektur & Prozess

| # | Idee | Beschreibung |
|---|------|-------------|
| Q-01 | **Multi-Repo-Tracing** | Requirements in Repo A, Implementierung in Repo B (Issue #88) |
| Q-02 | **Plugin-System** | Eigene Parser/Validatoren als Go-Plugins (Issue #90) |
| Q-03 | **Risikomatrix** | `[risk,id=RISK-001,fmea=...]` Blöcke mit FMEA-Verknüpfung (Issue #95) |
| Q-04 | **FMEA-Integration** | Failure Mode and Effects Analysis verknüpfen |
| Q-05 | **Safety-Goals** | ISO 26262 ASIL-Level auf Anforderungen |
| Q-06 | **Requirement-Versioning** | `version=2` tracen, Änderungshistorie visualisieren |
| Q-07 | **Baseline-Management** | Snapshot der Traceability zu einem Zeitpunkt einfrieren |
| Q-08 | **Review-Workflow** | `reviewed-by=` / `approved-by=` als Pflichtfelder erzwingen |
| Q-09 | **Traceability-Score** | Ein einziger Gesundheits-Score (0–100) für das Projekt |
| Q-10 | **SYS-Level Requirements** | Systemanforderungen → SW-Anforderungen ableiten (Issue #153) |
| Q-11 | **Stakeholder-View** | Welche Requirements interessieren welchen Stakeholder |
| Q-12 | **Test-Pyramide** | Unit/Integration/System-Tests kategorisieren und visualisieren |
| Q-13 | **Mutation-Testing** | Mutation-Score als Test-Qualitätsmetrik einbinden |
| Q-14 | **Requirement-Diff** | Änderungen zwischen Commits als Diff darstellen (Issue #83) |

---

## 🛠️ Developer Experience

| # | Idee | Beschreibung |
|---|------|-------------|
| D-01 | **`init --template aspice`** | Vorkonfigurierte Templates für ASPICE-Projekte |
| D-02 | **Watch-Modus Partial-Reload** | Nur geänderte Dateien neu laden statt vollständigem Rebuild |
| D-03 | **`req42-tracer check`** | Schneller Pre-Commit-Check (< 100ms) |
| D-04 | **Shell-Completion** | Bash/Zsh/Fish-Autovervollständigung |
| D-05 | **`req42-tracer explain REQ-001`** | Zeigt alle Traces für eine ID direkt im Terminal |
| D-06 | **`req42-tracer stats`** | Schnelle Statistik ohne vollen Report zu generieren |
| D-07 | **`req42-tracer lint`** | Nur Syntax/Konsistenz prüfen, kein Report |
| D-08 | **Interactive TUI** | Terminal-UI (bubbletea) für Navigation in der Traceability |
| D-09 | **`req42-tracer serve`** | Lokaler HTTP-Server mit Live-Reload statt Watch-Mode |
| D-10 | **Config-Wizard** | Interaktive `.req42.yaml`-Erstellung |
| D-11 | **Migration-Tool** | Alte Dokumente auf neues Format migrieren |
| D-12 | **`req42-tracer diff HEAD~1`** | Traceability-Änderungen zwischen Commits anzeigen |

---

## 🧪 Testing & Qualität

| # | Idee | Beschreibung |
|---|------|-------------|
| T-01 | **Property-Based Testing** | Fuzzing für den AsciiDoc-Parser |
| T-02 | **Golden-File-Tests** | Report-Output mit gespeicherten Snapshots vergleichen |
| T-03 | **Performance-Benchmarks** | Wie schnell ist der Parser bei 10.000 Anforderungen? |
| T-04 | **Load-Testing** | HTML-Report mit 1000+ Nodes |
| T-05 | **Accessibility** | WCAG-konformer HTML-Report (Screen Reader, Kontrast) |
| T-06 | **`req42-tracer test`** | Integrierter Test-Runner mit Coverage-Gate |

---

## 📦 Distribution & Deployment

| # | Idee | Beschreibung |
|---|------|-------------|
| X-01 | **Homebrew-Formula** | `brew install req42-tracer` |
| X-02 | **apt/rpm-Package** | Linux-Paketmanager-Support |
| X-03 | **Docker-Image** | `docker run req42-tracer trace` |
| X-04 | **WASM-Build** | req42-tracer im Browser laufen lassen |
| X-05 | **GitHub Action** | `uses: paul-fleischmann-com/req42-tracer-action@v1` |
| X-06 | **npm-Package** | JavaScript-Wrapper für Node.js-Projekte |
| X-07 | **pip-Package** | Python-Wrapper |

---

## 🌐 Standards & Compliance

| # | Idee | Beschreibung |
|---|------|-------------|
| S-01 | **ISO 25010-Mapping** | Software-Qualitätseigenschaften auf Requirements mappen |
| S-02 | **IEC 61508-Support** | Funktionale Sicherheit (SIL-Level) |
| S-03 | **DO-178C-Support** | Avionik-Software-Standard |
| S-04 | **MISRA-Verlinkung** | MISRA-Regeln mit Anforderungen verknüpfen |
| S-05 | **AUTOSAR-Integration** | AUTOSAR-Komponenten als Arch-Elemente |
| S-06 | **SysML/UML-Import** | Block Definition Diagrams als Arch-Elemente importieren |

---

*Letzte Aktualisierung: 2026-05-31 — Ideen ohne Priorisierung, dienen als Diskussionsbasis.*
