# Setup

Einmalige Einrichtung der Entwicklungsumgebung für req42-tracer.

## Voraussetzungen

- Go ≥ 1.21
- git
- [gh CLI](https://cli.github.com/) (GitHub CLI)

## 1. Repositories klonen

```bash
# Ins Workspace-Verzeichnis wechseln
cd /workspace

# req42-tracer
git clone https://github.com/paul-fleischmann-com/req42-tracer

# Claude Skill Library (auf gleicher Ebene)
git clone https://github.com/paul-fleischmann-com/claude-skill-library
```

## 2. Claude Skill Library installieren

```bash
# Installation (einmalig) — installiert Skills nach ~/.claude/skill-library
REPO_URL="https://github.com/paul-fleischmann-com/claude-skill-library.git"
TARGET="$HOME/.claude/skill-library"

git clone "$REPO_URL" "$TARGET"
mkdir -p "$HOME/.claude"
ln -sfn "$TARGET/skills" "$HOME/.claude/skills"
```

> **Hinweis:** `install.sh` im Repo verwendet SSH. Für HTTPS-Umgebungen die obigen Befehle direkt ausführen.

## 3. Skills aktuell halten

```bash
# Skills updaten (zieht neueste Version aus GitHub)
bash /workspace/claude-skill-library/update.sh
```

## 4. GitHub CLI einrichten

Beide Accounts einloggen:

```bash
# Reviewer-Account (Admin)
echo "TOKEN" | gh auth login --with-token

# Developer-Account (Write)
echo "TOKEN" | gh auth login --with-token

# Account wechseln
gh auth switch --user paulefl          # für Review & Merge
gh auth switch --user dev-paul-fleischmann  # für Implementierung
```

Siehe [`ROLES.md`](ROLES.md) für Details zur Rollentrennung.

## 5. Go-Abhängigkeiten

```bash
cd /workspace/req42-tracer
go mod download
go build ./...
go test ./...
```

## Verzeichnisstruktur nach Setup

```
/workspace/
  req42-tracer/          # dieses Repo
  claude-skill-library/  # Claude Skills (workspace-Klon)

~/.claude/
  skill-library/         # installierte Skills (geklont von GitHub)
  skills -> skill-library/skills  # Symlink
```
