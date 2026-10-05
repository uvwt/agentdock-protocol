package protocol

const ConnectionProtocolVersion = "2"

const (
	MessageNodeHello     = "node.hello"
	MessageNodeReady     = "node.ready"
	MessageNodeHeartbeat = "node.heartbeat"
	MessageToolInvoke    = "tool.invoke"
	MessageToolResult    = "tool.result"
	MessageToolError     = "tool.error"
	MessageToolCancel    = "tool.cancel"
)

const (
	OperationRuntimeRequest = "runtime.request"
	OperationContextLocal   = "context.local"
	OperationToolCall       = "tool.call"
	OperationResourceRead   = "resource.read"
	OperationArtifactRead   = "artifact.read"
)

const (
	// CapabilitiesNegotiationCapability marks a node that uses explicit Bridge feature negotiation.
	// Nodes without this marker are legacy v2 nodes and retain the historical optional-operation behavior.
	CapabilitiesNegotiationCapability = "bridge.capabilities.v1"
	// ContextLocalCapability advertises support for context.local.
	ContextLocalCapability = "bridge.context.local.v1"
	// RuntimeRequestCapability advertises support for runtime.request.
	RuntimeRequestCapability = "bridge.runtime.request.v1"
	// ResourceReadCapability advertises support for resource.read.
	ResourceReadCapability = "bridge.resource.read.v1"
	// ArtifactReadCapability is the frozen Bridge capability token for private Artifact reads.
	ArtifactReadCapability = "bridge.artifact.read.v1"
	// MaxArtifactChunkBytes bounds one raw Artifact chunk before JSON/base64 encoding on the Bridge.
	MaxArtifactChunkBytes = 512 << 10
)

const AgentDockUIResourcePrefix = "ui://agentdock/"

const (
	ContextUIResourceURI      = "ui://agentdock/context"
	WorkspaceUIResourceURI    = "ui://agentdock/workspace-context"
	TaskProgressUIResourceURI = "ui://agentdock/task-progress"
	FileChangeUIResourceURI   = "ui://agentdock/file-change"
	RecallUIResourceURI       = "ui://agentdock/recall"
	WorkflowUIResourceURI     = "ui://agentdock/workflow"
	DynamicMCPUIResourceURI   = "ui://agentdock/dynamic-mcp"
	ArtifactUIResourceURI     = "ui://agentdock/artifact"
	ACPStatusUIResourceURI    = "ui://agentdock/acp-status"
	ACPPromptUIResourceURI    = "ui://agentdock/acp-prompt"
	ImageUIResourceURI        = "ui://agentdock/view-image/v1.html"
)

const (
	ContextUIContract      = "agentdock.context.fleet.v1"
	WorkspaceUIContract    = "agentdock.workspace-context.v1"
	TaskProgressUIContract = "agentdock.task-progress.v1"
	FileChangeUIContract   = "agentdock.file-change.v1"
	RecallUIContract       = "agentdock.recall.v1"
	WorkflowUIContract     = "agentdock.workflow.v1"
	DynamicMCPUIContract   = "agentdock.dynamic-mcp.v1"
	ArtifactUIContract     = "agentdock.artifact.v1"
	ACPStatusUIContract    = "agentdock.acp-status.v1"
	ACPPromptUIContract    = "agentdock.acp-prompt.v1"
	ImageUIContract        = "agentdock.view-image.v1"
)

const MCPAppMIMEType = "text/html;profile=mcp-app"

// UIResourceContract returns the renderer contract bound to one AgentDock MCP App URI.
// URIs identify resources and remain stable; renderer compatibility evolves through the contract string.
func UIResourceContract(uri string) (string, bool) {
	switch uri {
	case ImageUIResourceURI:
		return ImageUIContract, true
	case ContextUIResourceURI:
		return ContextUIContract, true
	case WorkspaceUIResourceURI:
		return WorkspaceUIContract, true
	case TaskProgressUIResourceURI:
		return TaskProgressUIContract, true
	case FileChangeUIResourceURI:
		return FileChangeUIContract, true
	case RecallUIResourceURI:
		return RecallUIContract, true
	case WorkflowUIResourceURI:
		return WorkflowUIContract, true
	case DynamicMCPUIResourceURI:
		return DynamicMCPUIContract, true
	case ArtifactUIResourceURI:
		return ArtifactUIContract, true
	case ACPStatusUIResourceURI:
		return ACPStatusUIContract, true
	case ACPPromptUIResourceURI:
		return ACPPromptUIContract, true
	default:
		return "", false
	}
}
