# Review Workflow

Dieser Prozess gilt für alle Code- und Security-Reviews in diesem Repository.

## Code Review

Review starten mit `/code-review` (Effort: `low` / `medium` / `high` / `ultra`).

### Ablauf

1. `/code-review high` — Review ausführen
2. Alle **CONFIRMED** und **PLAUSIBLE** Findings als Inline-Kommentar im PR dokumentieren
3. Gefundene Bugs sofort fixen, committen und pushen
4. Inline-Kommentar an die geänderte Zeile setzen (nach dem Fix)
5. **REFUTED** Findings werden nicht kommentiert

### Inline-Kommentar setzen

```bash
gh api repos/OWNER/REPO/pulls/PR_NUMBER/comments \
  --method POST \
  --field body="**Bug (Code Review):** Beschreibung des Bugs und der Ursache." \
  --field commit_id="COMMIT_SHA" \
  --field path="path/to/file.go" \
  --field line=LINE_NUMBER \
  --field side="RIGHT"
```

---

## Security Review

Review starten mit `/security-review`.

### Schweregrade

| Symbol | Stufe    | Bedeutung                                      |
|--------|----------|------------------------------------------------|
| 🔴     | Critical | Sofort fixen — kein Merge ohne Fix             |
| 🟠     | High     | Fix im selben PR                               |
| 🟡     | Medium   | Fix in separatem Issue/PR                      |
| 🔵     | Low      | Dokumentieren, kein Blockierungsgrund          |

### Ablauf

1. `/security-review` — Review ausführen
2. Alle Findings als Inline-Kommentar im PR dokumentieren (Schweregrad angeben)
3. 🔴 Critical und 🟠 High sofort fixen vor dem Merge
4. 🟡 Medium und 🔵 Low als GitHub Issue anlegen

### Inline-Kommentar setzen

```bash
gh api repos/OWNER/REPO/pulls/PR_NUMBER/comments \
  --method POST \
  --field body="🔴 **Critical (Security Review):** Beschreibung des Findings." \
  --field commit_id="COMMIT_SHA" \
  --field path="path/to/file.go" \
  --field line=LINE_NUMBER \
  --field side="RIGHT"
```

---

## Traceability-Check

Bei jeder Implementierung sicherstellen, dass Requirements, Architektur und Tests korrekt verlinkt sind.

### Checkliste

Für jede neue Implementierung (Paket, Command, Feature) prüfen:

| # | Prüfpunkt | Befehl / Ort |
|---|-----------|-------------|
| 1 | **`[req,id=REQ-*]` vorhanden?** Ist eine Anforderung in `docs/requirements/req42.adoc` die das Feature beschreibt? | `grep "REQ-XYZ" docs/requirements/req42.adoc` |
| 2 | **`[arch,id=comp.*,impl=<pfad>,req=REQ-*]` vorhanden?** Hat das Arch-Element ein gültiges `impl=`-Feld das auf den echten Pfad zeigt? | `grep "comp.xyz" docs/arc42/arc42.adoc` |
| 3 | **`impl=`-Pfad existiert auf Disk?** Stimmt der `impl=`-Pfad mit der tatsächlichen Datei/Package überein? | `ls <impl-pfad>` |
| 4 | **`architecture.jsonc` aktualisiert?** Hat das JSONC-Element ebenfalls ein `impl=`-Feld? | `grep "backend_\|cli_" architecture.jsonc` |
| 5 | **`[test-spec,id=TS-PKG-NNN,req="REQ-*",aspice="SWE.5.BP3"]` vorhanden?** Hat jede Testfunktion einen korrekt formatierten Test-Spec-Kommentar? | `grep "test-spec" internal/<paket>/*_test.go` |
| 6 | **`req42-tracer validate` sauber?** Keine Validierungsfehler? | `go run ./cmd/req42-tracer validate` |
| 7 | **`req42-tracer gaps` sauber?** Keine MissingImplementation-Gaps (außer bekannte wie comp.lsp)? | `go run ./cmd/req42-tracer gaps` |

### Test-Spec Format

```go
// [test-spec,id=TS-PKG-NNN,req="REQ-PKG-NNN",aspice="SWE.5.BP3"]
// TestFunctionName beschreibt was getestet wird.
func TestFunctionName(t *testing.T) {
```

| Feld | Pflicht | Schema | Beispiel |
|------|---------|--------|---------|
| `id` | ✅ | `TS-<PAKET>-<NNN>` | `TS-LSP-001` |
| `req` | ✅ | Quoted Req-ID | `"REQ-LSP-001"` |
| `aspice` | ✅ | Quoted ASPICE-BP | `"SWE.5.BP3"` |

### Typische Fehler

| Fehler | Symptom in `gaps` | Fix |
|--------|-------------------|-----|
| `impl=` fehlt in arc42.adoc | `MissingImplementation: comp.xyz` | `impl=internal/xyz` ergänzen |
| `impl=`-Pfad existiert nicht | validate-Warning | Pfad korrigieren |
| Test-Spec ohne Quotes | Parser ignoriert Block | `req="REQ-X"` statt `req=REQ-X` |
| Test-Spec-ID falsches Schema | Kein Trace in Report | `TS-PKG-NNN` statt `spec.pkg.name` |
| JSONC `impl=` fehlt | Dokumentations-Inkonsistenz | `"impl": "internal/xyz"` ergänzen |

### Ablauf

1. Nach Implementierung: Checkliste oben Punkt für Punkt durchgehen
2. `req42-tracer validate && req42-tracer gaps` ausführen — beide müssen sauber sein
3. Ergebnis im PR-Review vermerken

---

## Bausteinsicht-Check

Bei jeder Implementierung prüfen, ob `architecture.jsonc` aktualisiert werden muss.

### Checkliste

1. **Neues CLI-Command** → neues `cli_<name>`-Element unter dem `cli`-Container?
2. **Neuer interner Service/Server** (z.B. HTTP-Server im Watch-Command) → eigenes Element oder Beschreibung anpassen?
3. **Neue externe Abhängigkeit** (Library, Protokoll, Datenquelle) → neues `External System`-Element?
4. **Neue View** nötig (z.B. wenn ein neuer Deployment-Kontext entsteht)?

### Ablauf

1. `architecture.jsonc` öffnen und mit der implementierten Funktionalität abgleichen
2. Fehlende Elemente ergänzen (ID-Schema: `<container>_<name>`, z.B. `cli_watch`)
3. Fehlende View-Einträge ergänzen
4. Ergebnis im PR-Review vermerken (auch wenn keine Änderung nötig war)

### Beispiel (Phase 2.4 — Watch Mode)

| Prüfpunkt | Ergebnis |
|---|---|
| `cli_watch`-Element vorhanden? | ✅ bereits in `architecture.jsonc` |
| HTTP-Server als separates Element? | ℹ️ Implementierungsdetail von `cli_watch` — kein eigenes Element nötig |
| Neue externe Abhängigkeit (fsnotify)? | ℹ️ Go-Library, kein externes System im C4-Sinne |
| View-Eintrag `Component — CLI` enthält `cli_watch`? | ✅ vorhanden |
| Projekt-Key korrekt (`req42-project`)? | ✅ nach Rename aktuell |

---

## Hinweise

- `commit_id` = aktueller HEAD-Commit des PR-Branches (`git rev-parse HEAD`)
- `side` = `RIGHT` für neue/geänderte Zeilen, `LEFT` für entfernte Zeilen
- `line` = Zeilennummer in der Datei (nicht im Diff)
- Kommentare landen direkt als Review-Threads im PR und sind für alle Reviewer sichtbar
