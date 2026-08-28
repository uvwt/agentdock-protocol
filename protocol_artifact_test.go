package protocol

import (
	"encoding/json"
	"testing"
)

func TestArtifactBridgeWireLiteralsStayStable(t *testing.T) {
	if ConnectionProtocolVersion != "2" {
		t.Fatalf("ConnectionProtocolVersion = %q, want 2", ConnectionProtocolVersion)
	}
	if OperationArtifactRead != "artifact.read" {
		t.Fatalf("OperationArtifactRead = %q", OperationArtifactRead)
	}
	if ArtifactReadCapability != "bridge.artifact.read.v1" {
		t.Fatalf("ArtifactReadCapability = %q", ArtifactReadCapability)
	}
	if MaxArtifactChunkBytes != 512<<10 {
		t.Fatalf("MaxArtifactChunkBytes = %d", MaxArtifactChunkBytes)
	}
}

func TestHelloAdditiveFieldsRemainWireCompatible(t *testing.T) {
	encoded := []byte(`{
		"type":"node.hello",
		"protocol_version":"2",
		"hello":{
			"device_id":"device_abcdefgh",
			"version":"0.8.0",
			"protocol_version":"2",
			"os":"darwin",
			"arch":"arm64",
			"capabilities":["read_file"],
			"bridge_capabilities":["bridge.artifact.read.v1"],
			"tool_contract_hash":"",
			"tools":[],
			"ui_resources":[],
			"future_additive_field":{"supported":true}
		}
	}`)

	var message Message
	if err := json.Unmarshal(encoded, &message); err != nil {
		t.Fatalf("additive Hello field must be ignored by older-compatible decoding: %v", err)
	}
	if message.Hello == nil || len(message.Hello.BridgeCapabilities) != 1 || message.Hello.BridgeCapabilities[0] != ArtifactReadCapability {
		t.Fatalf("decoded hello = %#v", message.Hello)
	}
}
