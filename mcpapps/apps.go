package mcpapps

import (
	_ "embed"
	"strings"
)

//go:embed app.html
var appHTMLTemplate string

//go:embed image.html
var imageHTML string

// ImageResultText explains that the image can be explicitly handed to the model from the Image card.
const ImageResultText = "Image returned. Use the Image card to provide it to the model when needed."

// HTML renders one shared MCP App document for a known AgentDock view.
// View and title are internal constants owned by AgentDock/NexusDock, not user input.
func HTML(view, title string) string {
	if view == "view_image" {
		return imageHTML
	}
	return strings.NewReplacer("{{VIEW}}", view, "{{TITLE}}", title).Replace(appHTMLTemplate)
}
