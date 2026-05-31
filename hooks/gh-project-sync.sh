#!/usr/bin/env bash
# hooks/gh-project-sync.sh
# Syncs GitHub Issues/PRs automatically to the Project board.
#
# Setup: Erweitere das gh-Token um den 'project' Scope:
#   https://github.com/settings/tokens
#   Dann: gh auth refresh -h github.com -s project
#
# Konfiguration: .github/project-config.json anlegen (siehe unten).
# Das Script legt die Datei beim ersten erfolgreichen Lauf automatisch an.

set -euo pipefail

REPO="paul-fleischmann-com/req42-tracer"
OWNER="paul-fleischmann-com"
CONFIG_FILE=".github/project-config.json"

# Lese Hook-Input
HOOK_INPUT="${CLAUDE_HOOK_INPUT:-{}}"
TOOL_CMD=$(echo "$HOOK_INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null || echo "")
TOOL_OUTPUT=$(echo "$HOOK_INPUT" | jq -r '.tool_response.output // ""' 2>/dev/null || echo "")

# --- Kein relevanter Befehl? Sofort beenden ---
if ! echo "$TOOL_CMD" | grep -qE 'gh\s+(issue\s+create|pr\s+create|pr\s+merge)'; then
    exit 0
fi

# --- Scope-Check: 'project' erforderlich ---
if ! gh auth status 2>&1 | grep -q "project"; then
    # Scope fehlt noch — silent skip, keine Fehlermeldung
    exit 0
fi

# --- Config laden oder bootstrap ---
load_config() {
    if [ ! -f "$CONFIG_FILE" ]; then
        # Ersten verfügbaren ProjectV2 des Repos ermitteln
        local project_data
        project_data=$(gh api graphql -f query='
          query($owner:String!, $repo:String!) {
            repository(owner:$owner, name:$repo) {
              projectsV2(first:1) {
                nodes { id number title
                  fields(first:20) {
                    nodes {
                      ... on ProjectV2SingleSelectField {
                        id name
                        options { id name }
                      }
                    }
                  }
                }
              }
            }
          }' -f owner="$OWNER" -f repo="$(basename $REPO)" 2>/dev/null) || return 1

        local project_id project_num status_field_id
        project_id=$(echo "$project_data" | jq -r '.data.repository.projectsV2.nodes[0].id // empty')
        project_num=$(echo "$project_data" | jq -r '.data.repository.projectsV2.nodes[0].number // empty')
        [ -z "$project_id" ] && return 1

        status_field_id=$(echo "$project_data" | jq -r '
          .data.repository.projectsV2.nodes[0].fields.nodes[]
          | select(.name == "Status") | .id // empty' | head -1)
        [ -z "$status_field_id" ] && return 1

        mkdir -p .github
        echo "$project_data" | jq --arg pid "$project_id" --arg pnum "$project_num" --arg sfid "$status_field_id" '{
          project_id: $pid,
          project_number: ($pnum | tonumber),
          status_field_id: $sfid,
          status_options: (
            .data.repository.projectsV2.nodes[0].fields.nodes[]
            | select(.name == "Status")
            | .options
            | map({(.name): .id})
            | add
          )
        }' > "$CONFIG_FILE"
        echo "[gh-project-sync] Config geschrieben: $CONFIG_FILE" >&2
    fi
    return 0
}

load_config || exit 0
[ ! -f "$CONFIG_FILE" ] && exit 0

PROJECT_ID=$(jq -r '.project_id' "$CONFIG_FILE")
STATUS_FIELD_ID=$(jq -r '.status_field_id' "$CONFIG_FILE")

get_status_option() {
    local name="$1"
    jq -r --arg n "$name" '.status_options[$n] // empty' "$CONFIG_FILE"
}

# --- Issue einer URL dem Projekt hinzufügen + Status setzen ---
add_issue_to_project() {
    local issue_url="$1"
    local status_name="$2"

    local issue_num
    issue_num=$(echo "$issue_url" | grep -oE '[0-9]+$') || return 0

    local node_id
    node_id=$(gh api "repos/$REPO/issues/$issue_num" --jq '.node_id' 2>/dev/null) || return 0

    local item_id
    item_id=$(gh api graphql -f query='
      mutation($proj:ID!, $nid:ID!) {
        addProjectV2ItemById(input:{projectId:$proj, contentId:$nid}) {
          item { id }
        }
      }' -f proj="$PROJECT_ID" -f nid="$node_id" \
      --jq '.data.addProjectV2ItemById.item.id' 2>/dev/null) || return 0

    local option_id
    option_id=$(get_status_option "$status_name")
    [ -z "$option_id" ] && return 0

    gh api graphql -f query='
      mutation($proj:ID!, $item:ID!, $field:ID!, $val:String!) {
        updateProjectV2ItemFieldValue(input:{
          projectId:$proj, itemId:$item, fieldId:$field,
          value:{singleSelectOptionId:$val}
        }) { projectV2Item { id } }
      }' -f proj="$PROJECT_ID" -f item="$item_id" \
         -f field="$STATUS_FIELD_ID" -f val="$option_id" > /dev/null 2>&1 || true

    echo "[gh-project-sync] Issue #$issue_num → Project ($status_name)" >&2
}

# --- Logik je Befehlstyp ---

if echo "$TOOL_CMD" | grep -qE 'gh\s+issue\s+create'; then
    # Issue-URL aus Output extrahieren (Format: https://github.com/.../issues/NNN)
    ISSUE_URL=$(echo "$TOOL_OUTPUT" | grep -oE 'https://github\.com/[^/]+/[^/]+/issues/[0-9]+' | head -1)
    [ -n "$ISSUE_URL" ] && add_issue_to_project "$ISSUE_URL" "Todo"

elif echo "$TOOL_CMD" | grep -qE 'gh\s+pr\s+create'; then
    # PR erstellt → verknüpfte Issues auf "In Progress"
    PR_URL=$(echo "$TOOL_OUTPUT" | grep -oE 'https://github\.com/[^/]+/[^/]+/pull/[0-9]+' | head -1)
    [ -z "$PR_URL" ] && exit 0
    PR_NUM=$(echo "$PR_URL" | grep -oE '[0-9]+$')
    # Issues aus PR-Body holen (Closes/Fixes/Refs #NNN)
    PR_BODY=$(gh api "repos/$REPO/pulls/$PR_NUM" --jq '.body // ""' 2>/dev/null) || exit 0
    echo "$PR_BODY" | grep -oE '(Closes|Fixes|Refs)\s+#[0-9]+' | grep -oE '#[0-9]+' | tr -d '#' | while read -r num; do
        add_issue_to_project "https://github.com/$REPO/issues/$num" "In Progress"
    done

elif echo "$TOOL_CMD" | grep -qE 'gh\s+pr\s+merge'; then
    # PR gemergt → verknüpfte Issues auf "Done"
    PR_NUM=$(echo "$TOOL_CMD" | grep -oE '\b[0-9]+\b' | head -1)
    [ -z "$PR_NUM" ] && exit 0
    PR_BODY=$(gh api "repos/$REPO/pulls/$PR_NUM" --jq '.body // ""' 2>/dev/null) || exit 0
    echo "$PR_BODY" | grep -oE '(Closes|Fixes)\s+#[0-9]+' | grep -oE '#[0-9]+' | tr -d '#' | while read -r num; do
        add_issue_to_project "https://github.com/$REPO/issues/$num" "Done"
    done
fi

exit 0
