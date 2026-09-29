package protocol

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"testing"
)

// goldenKeys decodes a golden JSON object and returns its keys at every nesting level.
func goldenKeys(t *testing.T, name string) []string {
	t.Helper()
	content, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]any
	if err := json.Unmarshal(content, &object); err != nil {
		t.Fatal(err)
	}
	var keys []string
	var walk func(prefix string, value map[string]any)
	walk = func(prefix string, value map[string]any) {
		for key, item := range value {
			keys = append(keys, prefix+key)
			if nested, ok := item.(map[string]any); ok {
				walk(prefix+key+".", nested)
			}
		}
	}
	walk("", object)
	sort.Strings(keys)
	return keys
}

func decodeStrict(t *testing.T, name string, target any) {
	t.Helper()
	content, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		t.Fatalf("%s does not match the Go contract: %v", name, err)
	}
}

func TestPairGoldenFixturesMatchWireShape(t *testing.T) {
	cases := map[string][]string{
		"pair_request.json":  {"code", "device_id", "name"},
		"pair_response.json": {"device_token", "node", "node.id"},
		"node_ready.json":    {"heartbeat_ms", "protocol_version", "public_url", "type"},
	}
	for name, want := range cases {
		if got := goldenKeys(t, name); strings.Join(got, ",") != strings.Join(want, ",") {
			t.Fatalf("%s keys = %v; want %v", name, got, want)
		}
	}

	var request PairRequest
	decodeStrict(t, "pair_request.json", &request)
	if !ValidDeviceID(request.DeviceID) || !ValidNodeName(request.Name) || request.Code == "" {
		t.Fatalf("golden pair request is not valid: %#v", request)
	}
	var response PairResponse
	decodeStrict(t, "pair_response.json", &response)
	if response.Node.ID == "" || response.DeviceToken == "" {
		t.Fatalf("golden pair response lacks node identity: %#v", response)
	}
	var ready Message
	decodeStrict(t, "node_ready.json", &ready)
	if ready.Type != MessageNodeReady || ready.ProtocolVersion != ConnectionProtocolVersion {
		t.Fatalf("golden ready message = %#v", ready)
	}
}

func TestPairWireShapeCarriesNoServerSideScope(t *testing.T) {
	// Pairing is bound to whatever the pairing code was issued for on the server. Neither
	// side may add account, tenant or workspace selectors to these messages.
	for _, name := range []string{"pair_request.json", "pair_response.json", "node_ready.json"} {
		for _, key := range goldenKeys(t, name) {
			lower := strings.ToLower(key)
			for _, forbidden := range []string{"tenant", "workspace", "account", "membership", "generation"} {
				if strings.Contains(lower, forbidden) {
					t.Fatalf("%s carries server-side scope field %q", name, key)
				}
			}
		}
	}
}

func TestPairResponseIgnoresAdditionalServerFields(t *testing.T) {
	// NexusDock currently also returns "ok" and "connect_path"; nodes must keep working
	// when a server adds or omits such fields.
	body := `{"ok":true,"node":{"id":"node_1","name":"workstation","enabled":true},"device_token":"t","connect_path":"/v1/nodes/connect"}`
	var response PairResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatal(err)
	}
	if response.Node.ID != "node_1" || response.DeviceToken != "t" {
		t.Fatalf("decoded pair response = %#v", response)
	}
}

func TestValidDeviceIDAcceptsGeneratedIDsAndRejectsMalformedOnes(t *testing.T) {
	valid := []string{
		"device_q7dH3xK9mZpL2vN8bR4tY6wA1sC5eF0g", // AgentDock: "device_" + base64url(24 random bytes)
		"device_abcdefgh",
		"a1234567",
		"A" + strings.Repeat("x", 127),
	}
	invalid := []string{"", "short", "_leading_underscore", "-leading-dash", "has space_1234", "slash/12345678", "A" + strings.Repeat("x", 128)}
	for _, id := range valid {
		if !ValidDeviceID(id) {
			t.Errorf("ValidDeviceID(%q) = false; want true", id)
		}
	}
	for _, id := range invalid {
		if ValidDeviceID(id) {
			t.Errorf("ValidDeviceID(%q) = true; want false", id)
		}
	}
}

func TestValidNodeNameCountsRunesNotBytes(t *testing.T) {
	if !ValidNodeName(strings.Repeat("节", MaxNodeNameRunes)) {
		t.Fatal("100 CJK runes should be a valid node name")
	}
	if ValidNodeName("") || ValidNodeName(strings.Repeat("a", MaxNodeNameRunes+1)) {
		t.Fatal("empty and over-long node names must be rejected")
	}
}
