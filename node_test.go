package protocol

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPairRequestEncodesExactlyTheClosedFieldSet(t *testing.T) {
	encoded, err := json.Marshal(PairRequest{Code: "pair-code", DeviceID: "device_abcdefgh", Name: "laptop"})
	if err != nil {
		t.Fatal(err)
	}
	const want = `{"code":"pair-code","device_id":"device_abcdefgh","name":"laptop"}`
	if string(encoded) != want {
		t.Fatalf("PairRequest JSON = %s; want %s", encoded, want)
	}
}

func TestPairResponseIgnoresServerSpecificFields(t *testing.T) {
	body := `{"ok":true,"node":{"id":"node_1","device_id":"device_abcdefgh","enabled":true},` +
		`"device_token":"secret","connect_path":"/v1/nodes/connect"}`
	var decoded PairResponse
	if err := json.Unmarshal([]byte(body), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Node.ID != "node_1" || decoded.DeviceToken != "secret" {
		t.Fatalf("decoded pair response = %#v", decoded)
	}
}

func TestValidDeviceIDAcceptsGeneratedIDsAndRejectsMalformedOnes(t *testing.T) {
	// AgentDock generates "device_" + base64url(24 random bytes).
	generated := "device_" + strings.Repeat("Ab-_9", 6) + "xy"
	for _, id := range []string{generated, "device_abcdefgh", "a1234567", "A" + strings.Repeat("b", 127)} {
		if !ValidDeviceID(id) {
			t.Fatalf("ValidDeviceID(%q) = false; want true", id)
		}
	}
	for _, id := range []string{"", "a123456", "-device_abcdefgh", "device abcdefgh", "device/abcdefgh", "A" + strings.Repeat("b", 128)} {
		if ValidDeviceID(id) {
			t.Fatalf("ValidDeviceID(%q) = true; want false", id)
		}
	}
}

func TestValidNodeNameCountsRunesNotBytes(t *testing.T) {
	if !ValidNodeName(strings.Repeat("节", MaxNodeNameRunes)) {
		t.Fatal("name with MaxNodeNameRunes multi-byte runes rejected")
	}
	if ValidNodeName("") || ValidNodeName(strings.Repeat("a", MaxNodeNameRunes+1)) {
		t.Fatal("empty or oversized node name accepted")
	}
}

func TestNodeOAuthCallbackPathMatchesNexusRoute(t *testing.T) {
	if got := NodeOAuthCallbackPath("node_1"); got != "/oauth/mcp/nodes/node_1/callback" {
		t.Fatalf("NodeOAuthCallbackPath = %q", got)
	}
}
