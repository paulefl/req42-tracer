#!/usr/bin/env bash
# generate-diagrams.sh — Export all Bausteinsicht views as Mermaid diagrams.
#
# Views and titles are read directly from architecture.jsonc via the
# bausteinsicht tool — no hardcoded lists needed.
#
# Usage:
#   ./scripts/generate-diagrams.sh                   # .mmd + .adoc + .md
#   ./scripts/generate-diagrams.sh --include adoc    # (default) AsciiDoc target
#   ./scripts/generate-diagrams.sh --include md      # Markdown target

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BAUSTEINSICHT="$REPO_ROOT/tools/bausteinsicht/bausteinsicht"
MODEL="$REPO_ROOT/architecture.jsonc"
OUT="$REPO_ROOT/docs/arc42/diagrams"
INCLUDE_FORMAT="adoc"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --include) INCLUDE_FORMAT="$2"; shift 2 ;;
    *) echo "Error: unknown argument '$1'" >&2; exit 1 ;;
  esac
done

if [[ "$INCLUDE_FORMAT" != "adoc" && "$INCLUDE_FORMAT" != "md" ]]; then
  echo "Error: --include must be 'adoc' or 'md'" >&2; exit 1
fi

mkdir -p "$OUT"

echo "Reading views from architecture.jsonc ..."
echo "Include format: $INCLUDE_FORMAT"
echo ""

# Export all views as JSON, save to temp file, then process with python3.
TMPJSON=$(mktemp)
trap 'rm -f "$TMPJSON"' EXIT

"$BAUSTEINSICHT" --format json --model "$MODEL" \
  export-diagram --diagram-format mermaid > "$TMPJSON"

python3 - "$OUT" "$TMPJSON" << 'PYEOF'
import json, sys, os, re

out_dir  = sys.argv[1]
json_path = sys.argv[2]

with open(json_path) as f:
    data = json.load(f)

count = 0
for item in data:
    view   = item["view"]
    source = item["source"].rstrip("\n")

    # Extract title from mermaid line "    title <TITLE>"
    title_match = re.search(r'^\s+title\s+(.+)$', source, re.MULTILINE)
    title = title_match.group(1).strip() if title_match else view

    mmd_path  = os.path.join(out_dir, f"{view}.mmd")
    adoc_path = os.path.join(out_dir, f"{view}.adoc")
    md_path   = os.path.join(out_dir, f"{view}.md")

    with open(mmd_path, "w") as f:
        f.write(source + "\n")

    with open(adoc_path, "w") as f:
        f.write(f".{title}\n[mermaid]\n....\n{source}\n....\n")

    with open(md_path, "w") as f:
        f.write(f"### {title}\n\n```mermaid\n{source}\n```\n")

    print(f"  ✓ {view}  ({title})")
    count += 1

print(f"\nDone — {count} views × 3 formats = {count * 3} files")
PYEOF
