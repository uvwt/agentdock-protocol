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
)

const (
	ContextUIResourceURI      = "ui://agentdock/context"
	TaskProgressUIResourceURI = "ui://agentdock/task-progress"
	FileChangeUIResourceURI   = "ui://agentdock/file-change"
	RecallUIResourceURI       = "ui://agentdock/recall"
	WorkflowUIResourceURI     = "ui://agentdock/workflow"
	DynamicMCPUIResourceURI   = "ui://agentdock/dynamic-mcp"
	ArtifactUIResourceURI     = "ui://agentdock/artifact"
	ACPStatusUIResourceURI    = "ui://agentdock/acp-status"
)

const (
	ContextUIContract      = "agentdock.context.fleet.v1"
	TaskProgressUIContract = "agentdock.task-progress.v1"
	FileChangeUIContract   = "agentdock.file-change.v1"
	RecallUIContract       = "agentdock.recall.v1"
	WorkflowUIContract     = "agentdock.workflow.v1"
	DynamicMCPUIContract   = "agentdock.dynamic-mcp.v1"
	ArtifactUIContract     = "agentdock.artifact.v1"
	ACPStatusUIContract    = "agentdock.acp-status.v1"
)

const MCPAppMIMEType = "text/html;profile=mcp-app"

// UIResourceContract returns the renderer contract bound to one AgentDock MCP App URI.
// URIs identify resources and remain stable; renderer compatibility evolves through the contract string.
func UIResourceContract(uri string) (string, bool) {
	switch uri {
	case ContextUIResourceURI:
		return ContextUIContract, true
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
	default:
		return "", false
	}
}
