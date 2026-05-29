package templates

import (
	"embed"
)

// FS embeds all template files for the `init` command.
// These are used as defaults and can be overridden by templates in the repo.
//
//go:embed req42.adoc arc42.adoc architecture.jsonc .req42.yaml .gitignore
var FS embed.FS
