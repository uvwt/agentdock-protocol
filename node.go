package protocol

import (
	"regexp"
	"unicode/utf8"
)

// Node endpoints are relative to the Nexus endpoint a node was configured with.
//
// NodePairPath accepts a POST with PairRequest. It returns 201 Created with PairResponse,
// 400 for a malformed body, device_id or name, 401 when the pairing code is unknown,
// expired or already used, and 409 when device_id is already paired. Only a successful
// pairing consumes the pairing code.
//
// NodeConnectPath is a WebSocket upgrade authenticated by "Authorization: Bearer <device_token>".
// An invalid or revoked token is rejected with 401 before the upgrade. The node sends
// MessageNodeHello first; Hello.DeviceID must equal the DeviceID sent when pairing.
// Nexus then replies with MessageNodeReady, whose PublicURL is an HTTPS origin without path.
const (
	NodePairPath    = "/v1/nodes/pair"
	NodeConnectPath = "/v1/nodes/connect"
)

// MaxNodeNameRunes bounds the human-readable node name sent during pairing.
const MaxNodeNameRunes = 100

var deviceIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{7,127}$`)

// PairRequest is the body a node sends to NodePairPath with a one-time pairing code.
//
// The field set is closed: servers reject unknown fields, so a new field needs a protocol
// version rather than an optional extension. Nodes generate a fresh random DeviceID for every
// pairing, so servers reject an already paired DeviceID instead of replacing that node.
type PairRequest struct {
	Code     string `json:"code"`
	DeviceID string `json:"device_id"`
	Name     string `json:"name"`
}

// PairResponse is the 201 Created body of NodePairPath.
// Servers may add fields; nodes must ignore fields they do not know.
type PairResponse struct {
	Node        PairedNode `json:"node"`
	DeviceToken string     `json:"device_token"`
}

// PairedNode is the part of the paired node record that nodes rely on.
type PairedNode struct {
	ID string `json:"id"`
}

// ValidDeviceID reports whether id matches the device_id wire format.
func ValidDeviceID(id string) bool {
	return deviceIDPattern.MatchString(id)
}

// ValidNodeName reports whether an already trimmed node name fits the pairing contract.
func ValidNodeName(name string) bool {
	count := utf8.RuneCountInString(name)
	return count >= 1 && count <= MaxNodeNameRunes
}

// NodeOAuthCallbackPath is the path under MessageNodeReady.PublicURL that a node registers as
// its OAuth redirect URL for remote MCP servers. Nexus forwards the callback to that node, and
// the node validates its own OAuth state. nodeID must not contain "/".
func NodeOAuthCallbackPath(nodeID string) string {
	return "/oauth/mcp/nodes/" + nodeID + "/callback"
}
