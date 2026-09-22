package mcpapps

import (
	_ "embed"
	"strings"
)

//go:embed app.html
var appHTMLTemplate string

//go:embed image.html
var imageHTML string

// ImageResultText explains the automatic host handoff without claiming the current model turn saw pixels.
const ImageResultText = "Image returned. The image component automatically provides the pixels to the host model context when supported."

// HTML renders one shared MCP App document for a known AgentDock view.
// View and title are internal constants owned by AgentDock/NexusDock, not user input.
func HTML(view, title string) string {
	if view == "view_image" {
		return imageHTML
	}
	return strings.NewReplacer("{{VIEW}}", view, "{{TITLE}}", title).Replace(appHTMLTemplate)
}
