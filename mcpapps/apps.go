package mcpapps

import (
	_ "embed"
	"strings"
)

//go:embed app.html
var appHTMLTemplate string

// HTML renders one shared MCP App document for a known AgentDock view.
// View and title are internal constants owned by AgentDock/NexusDock, not user input.
func HTML(view, title string) string {
	return strings.NewReplacer("{{VIEW}}", view, "{{TITLE}}", title).Replace(appHTMLTemplate)
}
