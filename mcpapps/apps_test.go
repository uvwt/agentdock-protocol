package mcpapps

import (
	"strings"
	"testing"
)

func TestHTMLRendersSharedMCPAppTemplate(t *testing.T) {
	html := HTML("recall", "Recall")
	for _, marker := range []string{
		`<title>Recall</title>`,
		`expectedView="recall"`,
		`rpcRequest("ui/initialize"`,
		`protocolVersion:"2026-01-26"`,
		`rpcNotify("ui/notifications/initialized"`,
		`message.method==="ui/notifications/tool-result"`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("shared MCP App missing marker %q", marker)
		}
	}
	for _, forbidden := range []string{"{{VIEW}}", "{{TITLE}}", "window.openai", "openai/widget", "openai/outputTemplate"} {
		if strings.Contains(html, forbidden) {
			t.Fatalf("shared MCP App contains forbidden marker %q", forbidden)
		}
	}
}

func TestHTMLUsesStandardHostContextThemeAndSemanticTokens(t *testing.T) {
	html := HTML("workflow", "Workflow")
	for _, marker := range []string{
		`:root{color-scheme:light dark;`,
		`:root[data-theme="dark"]{color-scheme:dark;`,
		`@media(prefers-color-scheme:dark)`,
		`--ad-text-primary:`,
		`var(--ad-text-primary)`,
		`const initialized=await rpcRequest("ui/initialize"`,
		`applyHostContext(initialized&&initialized.hostContext)`,
		`message.method==="ui/notifications/host-context-changed"`,
		`applyHostContext(message.params)`,
		`context.styles.variables`,
		`window.matchMedia("(prefers-color-scheme: dark)")`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("shared MCP App missing theme marker %q", marker)
		}
	}
	if strings.Contains(html, `:root{color-scheme:light;font:`) {
		t.Fatal("shared MCP App still forces a light-only color scheme")
	}
}
