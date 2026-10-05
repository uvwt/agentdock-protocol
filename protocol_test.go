package protocol

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestUIResourceContractsAreExplicit(t *testing.T) {
	cases := map[string]string{
		ContextUIResourceURI:      ContextUIContract,
		WorkspaceUIResourceURI:    WorkspaceUIContract,
		TaskProgressUIResourceURI: TaskProgressUIContract,
		FileChangeUIResourceURI:   FileChangeUIContract,
		RecallUIResourceURI:       RecallUIContract,
		WorkflowUIResourceURI:     WorkflowUIContract,
		DynamicMCPUIResourceURI:   DynamicMCPUIContract,
		ArtifactUIResourceURI:     ArtifactUIContract,
		ACPStatusUIResourceURI:    ACPStatusUIContract,
		ACPPromptUIResourceURI:    ACPPromptUIContract,
	}
	for uri, want := range cases {
		got, ok := UIResourceContract(uri)
		if !ok || got != want {
			t.Fatalf("UIResourceContract(%q) = %q, %v; want %q, true", uri, got, ok, want)
		}
	}
	if _, ok := UIResourceContract("ui://agentdock/unknown"); ok {
		t.Fatal("unknown AgentDock UI resource unexpectedly has a contract")
	}
}

func TestHelloRoundTripKeepsUIResourcesSeparateFromToolMeta(t *testing.T) {
	original := Message{
		Type:            MessageNodeHello,
		ProtocolVersion: ConnectionProtocolVersion,
		Hello: &Hello{
			DeviceID:           "device_abcdefgh",
			ProtocolVersion:    ConnectionProtocolVersion,
			Capabilities:       []string{"read_file"},
			BridgeCapabilities: []string{ArtifactReadCapability},
			Tools: []ToolDescriptor{{
				Name:        "workflow_template_manage",
				InputSchema: map[string]any{"type": "object"},
			}},
			UIResources: []UIResourceCapability{{
				URI: WorkflowUIResourceURI, Contract: WorkflowUIContract, MIMEType: MCPAppMIMEType,
			}},
		},
	}
	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Message
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ProtocolVersion != ConnectionProtocolVersion || decoded.Hello == nil {
		t.Fatalf("decoded message = %#v", decoded)
	}
	if len(decoded.Hello.UIResources) != 1 || decoded.Hello.UIResources[0].URI != WorkflowUIResourceURI {
		t.Fatalf("decoded ui_resources = %#v", decoded.Hello.UIResources)
	}
	if len(decoded.Hello.Capabilities) != 1 || decoded.Hello.Capabilities[0] != "read_file" {
		t.Fatalf("decoded capabilities = %#v", decoded.Hello.Capabilities)
	}
	if len(decoded.Hello.BridgeCapabilities) != 1 || decoded.Hello.BridgeCapabilities[0] != ArtifactReadCapability {
		t.Fatalf("decoded bridge_capabilities = %#v", decoded.Hello.BridgeCapabilities)
	}
	if decoded.Hello.Tools[0].Meta != nil {
		t.Fatalf("workflow tool unexpectedly carries static UI meta: %#v", decoded.Hello.Tools[0].Meta)
	}
}

func TestHelloCapabilityNegotiationFixtures(t *testing.T) {
	t.Run("legacy v2 hello remains valid without negotiation marker", func(t *testing.T) {
		message := readHelloFixture(t, "testdata/node_hello_legacy_v2.json")
		if message.ProtocolVersion != ConnectionProtocolVersion || message.Hello == nil {
			t.Fatalf("decoded legacy hello = %#v", message)
		}
		if len(message.Hello.BridgeCapabilities) != 1 || message.Hello.BridgeCapabilities[0] != ArtifactReadCapability {
			t.Fatalf("legacy bridge_capabilities = %#v", message.Hello.BridgeCapabilities)
		}
		for _, capability := range message.Hello.BridgeCapabilities {
			if capability == CapabilitiesNegotiationCapability {
				t.Fatal("legacy fixture unexpectedly opts into explicit negotiation")
			}
		}
	})

	t.Run("explicit v2 hello carries the full desktop capability set", func(t *testing.T) {
		message := readHelloFixture(t, "testdata/node_hello_explicit_capabilities_v2.json")
		if message.ProtocolVersion != ConnectionProtocolVersion || message.Hello == nil {
			t.Fatalf("decoded explicit hello = %#v", message)
		}
		want := []string{
			CapabilitiesNegotiationCapability,
			ContextLocalCapability,
			RuntimeRequestCapability,
			ResourceReadCapability,
			ArtifactReadCapability,
		}
		if !equalStrings(message.Hello.BridgeCapabilities, want) {
			t.Fatalf("explicit bridge_capabilities = %#v, want %#v", message.Hello.BridgeCapabilities, want)
		}
	})
}

func readHelloFixture(t *testing.T, path string) Message {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var message Message
	if err := json.Unmarshal(raw, &message); err != nil {
		t.Fatal(err)
	}
	return message
}

func equalStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestMessageTraceContextRoundTripIsBackwardCompatible(t *testing.T) {
	original := Message{
		Type:        MessageToolInvoke,
		RequestID:   "req_trace",
		Traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		Tracestate:  "rojo=00f067aa0ba902b7",
		Operation:   OperationToolCall,
		Arguments:   json.RawMessage(`{"tool":"read_file","arguments":{}}`),
	}
	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Message
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Traceparent != original.Traceparent || decoded.Tracestate != original.Tracestate {
		t.Fatalf("trace context round trip = %q / %q", decoded.Traceparent, decoded.Tracestate)
	}
	if decoded.Operation != OperationToolCall || decoded.RequestID != original.RequestID {
		t.Fatalf("existing envelope fields changed: %#v", decoded)
	}
}

func TestMessageTraceContextFieldsRemainOptional(t *testing.T) {
	encoded, err := json.Marshal(Message{
		Type:      MessageToolInvoke,
		RequestID: "req_legacy",
		Operation: OperationToolCall,
		Arguments: json.RawMessage(`{"tool":"read_file","arguments":{}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(encoded, []byte("traceparent")) || bytes.Contains(encoded, []byte("tracestate")) {
		t.Fatalf("optional trace fields leaked into legacy message: %s", encoded)
	}

	legacy := []byte(`{"type":"tool.invoke","request_id":"req_old","operation":"tool.call","arguments":{"tool":"read_file","arguments":{}}}`)
	var decoded Message
	if err := json.Unmarshal(legacy, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Traceparent != "" || decoded.Tracestate != "" {
		t.Fatalf("legacy message unexpectedly produced trace context: %#v", decoded)
	}
	if decoded.RequestID != "req_old" || decoded.Operation != OperationToolCall {
		t.Fatalf("legacy envelope changed: %#v", decoded)
	}
}

func TestCurrentMessageFieldsAreIgnoredByLegacyDecoder(t *testing.T) {
	type legacyMessage struct {
		Type      string          `json:"type"`
		RequestID string          `json:"request_id,omitempty"`
		Operation string          `json:"operation,omitempty"`
		Arguments json.RawMessage `json:"arguments,omitempty"`
	}
	encoded, err := json.Marshal(Message{
		Type:        MessageToolInvoke,
		RequestID:   "req_new",
		Traceparent: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
		Tracestate:  "rojo=00f067aa0ba902b7",
		Operation:   OperationToolCall,
		Arguments:   json.RawMessage(`{"tool":"read_file","arguments":{}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	var decoded legacyMessage
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Type != MessageToolInvoke || decoded.RequestID != "req_new" || decoded.Operation != OperationToolCall {
		t.Fatalf("legacy decoder lost existing fields: %#v", decoded)
	}
}
