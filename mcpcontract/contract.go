package mcpcontract

const (
	ToolAgentDockContext       = "agentdock_context"
	ToolRecallBootstrap        = "recall_bootstrap"
	ToolRecallSearch           = "recall_search"
	ToolRecallRead             = "recall_read"
	ToolRecallWrite            = "recall_write"
	ToolRecallMaintain         = "recall_maintain"
	ToolPrivateNoteManage      = "private_note_manage"
	ToolWorkflowTemplateManage = "workflow_template_manage"
)

var toolNames = []string{
	ToolAgentDockContext,
	ToolRecallBootstrap,
	ToolRecallSearch,
	ToolRecallRead,
	ToolRecallWrite,
	ToolRecallMaintain,
	ToolPrivateNoteManage,
	ToolWorkflowTemplateManage,
}

// ToolNames returns the canonical model-facing tools shared by AgentDock and NexusDock.
func ToolNames() []string { return append([]string(nil), toolNames...) }

func IsCanonicalTool(name string) bool {
	for _, candidate := range toolNames {
		if candidate == name {
			return true
		}
	}
	return false
}

type Annotations struct {
	ReadOnlyHint    bool
	DestructiveHint *bool
	IdempotentHint  *bool
	OpenWorldHint   *bool
}

func AnnotationContract(name string) (Annotations, bool) {
	readOnly := false
	destructive := true
	switch name {
	case ToolAgentDockContext, ToolRecallBootstrap, ToolRecallSearch, ToolRecallRead:
		readOnly = true
		destructive = false
	case ToolRecallWrite, ToolRecallMaintain, ToolPrivateNoteManage, ToolWorkflowTemplateManage:
	default:
		return Annotations{}, false
	}
	openWorld := false
	return Annotations{
		ReadOnlyHint:    readOnly,
		DestructiveHint: boolPtr(destructive),
		OpenWorldHint:   boolPtr(openWorld),
	}, true
}

func boolPtr(value bool) *bool { return &value }
