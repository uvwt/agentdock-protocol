package mcpapps

import (
	protocol "github.com/uvwt/agentdock-protocol"
	"strings"
	"testing"
)

func TestImageViewIsSharedAndHasStableContract(t *testing.T) {
	contract, ok := protocol.UIResourceContract("ui://agentdock/view-image/v1.html")
	if !ok || contract != "agentdock.view-image.v1" {
		t.Fatal("missing shared image contract")
	}
	html := HTML("view_image", "Image")
	for _, marker := range []string{"imageIds", "uploadFile", "ui/notifications/tool-result", "image-init"} {
		if !strings.Contains(html, marker) {
			t.Fatalf("missing image behavior: %s", marker)
		}
	}
	if strings.Contains(html, "nexusdock-image") {
		t.Fatal("shared renderer is tied to NexusDock")
	}
	if strings.Contains(HTML("agentdock_context", "Context"), "image-init") {
		t.Fatal("changed another view")
	}
}
