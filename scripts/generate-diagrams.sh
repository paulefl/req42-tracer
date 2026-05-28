#!/usr/bin/env bash
# generate-diagrams.sh — Export Bausteinsicht views as Mermaid diagrams.
#
# Usage:
#   ./scripts/generate-diagrams.sh           # generates .mmd + .adoc + .md
#   ./scripts/generate-diagrams.sh --include adoc   # default: arc42 includes .adoc files
#   ./scripts/generate-diagrams.sh --include md     # arc42 includes .md files
#
# Output: docs/arc42/diagrams/<view>.{mmd,adoc,md}

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BAUSTEINSICHT="$REPO_ROOT/tools/bausteinsicht/bausteinsicht"
MODEL="$REPO_ROOT/architecture.jsonc"
OUT="$REPO_ROOT/docs/arc42/diagrams"
INCLUDE_FORMAT="adoc"

# Parse args
while [[ $# -gt 0 ]]; do
  case "$1" in
    --include) INCLUDE_FORMAT="$2"; shift 2 ;;
    *) echo "Unknown argument: $1"; exit 1 ;;
  esac
done

if [[ "$INCLUDE_FORMAT" != "adoc" && "$INCLUDE_FORMAT" != "md" ]]; then
  echo "Error: --include must be 'adoc' or 'md'" >&2
  exit 1
fi

mkdir -p "$OUT"

VIEWS=(context containers backend backend_level2 frontend cli)

VIEW_TITLES=(
  "System Context"
  "Container View"
  "Component — Backend"
  "Component — Backend Level 2"
  "Component — Frontend"
  "Component — CLI"
)

echo "Generating diagrams → $OUT"
echo "Include format: $INCLUDE_FORMAT"
echo ""

for i in "${!VIEWS[@]}"; do
  VIEW="${VIEWS[$i]}"
  TITLE="${VIEW_TITLES[$i]}"

  # 1. Raw Mermaid source
  MMD_FILE="$OUT/${VIEW}.mmd"
  "$BAUSTEINSICHT" export-diagram \
    --model "$MODEL" \
    --diagram-format mermaid \
    --view "$VIEW" > "$MMD_FILE"
  echo "  ✓ ${VIEW}.mmd"

  # 2. AsciiDoc wrapper
  ADOC_FILE="$OUT/${VIEW}.adoc"
  {
    echo ".${TITLE}"
    echo "[mermaid]"
    echo "...."
    cat "$MMD_FILE"
    echo "...."
  } > "$ADOC_FILE"
  echo "  ✓ ${VIEW}.adoc"

  # 3. Markdown wrapper
  MD_FILE="$OUT/${VIEW}.md"
  {
    echo "### ${TITLE}"
    echo ""
    echo '```mermaid'
    cat "$MMD_FILE"
    echo '```'
  } > "$MD_FILE"
  echo "  ✓ ${VIEW}.md"
done

echo ""
echo "Done — ${#VIEWS[@]} views × 3 formats = $((${#VIEWS[@]} * 3)) files"
echo ""
echo "To include in arc42.adoc: --include $INCLUDE_FORMAT"
