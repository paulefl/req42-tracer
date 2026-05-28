# Project Roles

Dieses Dokument definiert die Rollen und Verantwortlichkeiten im req42-tracer Projekt.

## Rollen

### Developer — `dev-paul-fleischmann`

**Verantwortung:** Implementierung

- Feature-Branches erstellen und implementieren
- Commits auf Feature-Branches pushen
- Pull Requests öffnen
- Code Review Findings fixen und committen
- Tests schreiben (gemäß [`TESTS.md`](TESTS.md))
- CI-Fehler beheben

**Branch-Konvention:**
```
git checkout -b <issue-#>-kurzer-name
# Commits als dev-paul-fleischmann
```

**PR erstellen:**
```bash
gh pr create --assignee dev-paul-fleischmann --reviewer paulefl
```

---

### Reviewer — `paulefl`

**Verantwortung:** Code- und Security-Review, Merge-Entscheidung

- Pull Requests reviewen (Code Review + Security Review gemäß [`REVIEW.md`](REVIEW.md))
- Review-Findings als Inline-Kommentare im PR dokumentieren
- PRs approven oder Änderungen anfordern
- Feature-Branches in `master` mergen
- Releases taggen

---

## Workflow

```
dev-paul-fleischmann          paulefl
        │                        │
        │  feature branch        │
        ├──────────────────>     │
        │  implement + test      │
        │  /code-review          │
        │  /security-review      │
        │                        │
        │  open PR               │
        ├──────────────────────> │
        │                        │  review
        │                        │  inline comments
        │  fix findings  <───────┤
        ├──────────────────────> │
        │                        │  approve + merge
        │  <─────────────────────┤
```

## GitHub Konfiguration

| Setting | Wert |
|---|---|
| Default branch | `master` |
| Branch protection | PR required, 1 approval (`paulefl`) |
| Implementierung | `dev-paul-fleischmann` |
| Review & Merge | `paulefl` |
