package mcpapps

import (
	_ "embed"
	"strings"
)

//go:embed app.html
var appHTMLTemplate string

//go:embed image.html
var imageHTML string

// ImageResultText explains the optional, user-initiated follow-up without claiming the model saw pixels.
const ImageResultText = "图片已返回。若客户端未将图片直接提供给模型，可在图片组件中将它附加到下一轮对话。"

// HTML renders one shared MCP App document for a known AgentDock view.
// View and title are internal constants owned by AgentDock/NexusDock, not user input.
func HTML(view, title string) string {
	if view == "view_image" {
		return imageHTML
	}
	return strings.NewReplacer("{{VIEW}}", view, "{{TITLE}}", title).Replace(appHTMLTemplate)
}
