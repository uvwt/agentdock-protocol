package mcpcontract

// OutputSchema returns a fresh canonical output schema for tools whose result shape
// is shared by direct AgentDock and central NexusDock entrypoints.
func OutputSchema(name string) (map[string]any, bool) {
	props := map[string]any{}
	switch name {
	case ToolRecallBootstrap:
		props["recall_endpoint"] = stringProperty("NexusDock Recall endpoint.")
		props["project"] = stringProperty("Backend-selected Recall context; not an input selector for the model.")
		props["sections"] = arrayObjects("Packed Recall sections. Raw Markdown is returned only when include_raw=true.")
		props["count"] = integerProperty("Section count.")
		props["bytes"] = integerProperty("Combined bytes.")
	case ToolRecallSearch:
		props["recall_endpoint"] = stringProperty("NexusDock Recall endpoint.")
		props["recall_kind"] = stringProperty("Search kind used.")
		props["query"] = stringProperty("Search query.")
		props["recall_store"] = stringProperty("Recall store name.")
		props["results"] = map[string]any{
			"type": "array", "description": "Recall search results with source identity fields.",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"frontmatter":    objectProperty("Recall frontmatter metadata."),
					"matched_fields": arrayStrings("Matched fields."),
					"matched_terms":  arrayStrings("Matched terms."),
					"path":           stringProperty("Recall-relative document path."),
					"snippet":        stringProperty("Matched content snippet."),
					"id":             stringProperty("Stable Recall document id."),
					"title":          stringProperty("Human-readable document title."),
					"url":            stringProperty("Absolute source URL."),
				},
				"required": []string{"id", "title", "url"}, "additionalProperties": true,
			},
		}
		props["count"] = integerProperty("Search result count.")
	case ToolRecallRead:
		props["recall_endpoint"] = stringProperty("NexusDock Recall endpoint.")
		props["recall"] = objectProperty("NexusDock Recall document. Raw Markdown is returned only when include_raw=true.")
	case ToolRecallWrite:
		props["recall_endpoint"] = stringProperty("NexusDock Recall endpoint.")
		props["recall_target"] = stringProperty("Recall target used.")
		props["recall_action"] = stringProperty("Recall action used.")
		props["recall"] = objectProperty("NexusDock Recall document returned when a write occurs.")
		props["card"] = objectProperty("Normalized card candidate or written card when target=card.")
		props["warnings"] = arrayObjects("Review warnings before writing.")
		props["capture_plan"] = objectProperty("Reviewable write plan for card captures.")
		props["similar_results"] = arrayObjects("Similar existing card search results.")
		props["path"] = stringProperty("NexusDock Recall-relative path.")
		props["changed"] = booleanProperty("Whether the proposed edit changes content.")
		props["dry_run"] = booleanProperty("Whether the operation only previewed changes.")
		props["confirmed"] = booleanProperty("Whether write confirmation was supplied.")
		props["written"] = booleanProperty("Whether the entry was written.")
		props["diff"] = stringProperty("Unified diff preview.")
		props["updates"] = arrayObjects("Fact update results.")
	case ToolRecallMaintain:
		props["recall_endpoint"] = stringProperty("NexusDock Recall endpoint.")
		props["recall_action"] = stringProperty("Maintenance action performed.")
		props["entries"] = arrayObjects("NexusDock Recall entries for action=list.")
		props["count"] = integerProperty("Entry count where applicable.")
		props["terms"] = arrayStrings("Terms used for action=lint.")
		props["files_scanned"] = integerProperty("Files scanned for action=lint.")
		props["finding_count"] = integerProperty("Finding count for action=lint.")
		props["findings"] = arrayObjects("Lint findings.")
	case ToolPrivateNoteManage:
		props["root"] = stringProperty("NexusDock private notes root path.")
		props["private_note_store"] = stringProperty("Private note store name, fixed to NexusDock Private Notes.")
		props["action"] = stringProperty("Action performed.")
		props["query"] = stringProperty("Metadata-only query for action=search.")
		props["results"] = arrayObjects("Metadata-only private-note search results; never plaintext snippets.")
		props["metadata_only"] = booleanProperty("Whether search was restricted to safe metadata.")
		props["path"] = stringProperty("Plain note path for read/write/delete results.")
		props["encrypted_path"] = stringProperty("Age encrypted backup path.")
		props["content"] = stringProperty("Plaintext content returned only by explicit action=read.")
		props["truncated"] = booleanProperty("Whether returned content was truncated.")
		props["contains_secret"] = booleanProperty("Whether the note content is marked as containing secrets.")
		props["written"] = booleanProperty("Whether plaintext was written.")
		props["encrypted"] = booleanProperty("Whether encrypted backup was written.")
		props["deleted_plaintext"] = booleanProperty("Whether plaintext was deleted.")
		props["deleted_encrypted"] = booleanProperty("Whether encrypted backup was deleted.")
		props["notes"] = arrayObjects("Metadata-only private note summaries for status/list.")
		props["count"] = integerProperty("Result or note count.")
		props["notes_count"] = integerProperty("Private note count for status checks.")
		props["encrypted_count"] = integerProperty("Encrypted backup count for maintenance actions.")
		props["recipient"] = stringProperty("Age public recipient generated or used.")
		props["identity_created"] = booleanProperty("Whether a new local age identity was created.")
		props["algorithm"] = stringProperty("Encryption algorithm.")
		props["missing_encrypted"] = arrayStrings("Missing encrypted backup paths.")
		props["encrypted_backup_ok"] = booleanProperty("Whether every private note has its required encrypted backup.")
		props["plaintext_git_ignored"] = booleanProperty("Whether private note plaintext is Git-ignored.")
		props["keys_git_ignored"] = booleanProperty("Whether private note keys are Git-ignored.")
	case ToolWorkflowTemplateManage:
		props["action"] = stringProperty("Completed workflow template action.")
		props["template"] = objectProperty("Full workflow template returned by get.")
		props["templates"] = arrayObjects("Compact summaries from list or full active templates from get_many.")
		props["composition_required"] = booleanProperty("Whether the returned templates must be combined by the model before task creation.")
		props["next_required_action"] = stringProperty("Required model action after get_many.")
		props["template_id"] = stringProperty("Workflow template id returned by publish or retire.")
		props["template_summary"] = objectProperty("Compact workflow template summary returned by publish, retire, and list items.")
		props["count"] = integerProperty("Returned item count.")
		props["workflow_dir"] = stringProperty("Workflow template registry directory.")
		props["candidates"] = arrayObjects("Matched workflow template candidates with scores and reasons.")
		props["vector_search_enabled"] = booleanProperty("Whether optional embedding-backed template vector search is enabled for match.")
		props["vector_index_status"] = stringProperty("Template vector index status: disabled, ready, or degraded.")
		props["vector_index_items"] = integerProperty("Number of persisted template vectors for the current embedding model.")
		props["vector_index_available"] = booleanProperty("Whether workflow vector index content is available for export.")
		props["content"] = stringProperty("Raw workflow vector index JSON returned by vector_index.")
		props["embedding_model"] = stringProperty("Embedding model configured for template vector search.")
		props["recommended"] = stringProperty("Template recommendation: use_template, consider_template, or plain_task.")
		props["recommendation_reason"] = stringProperty("Reason for recommendation.")
		props["best_candidate_score"] = integerProperty("Highest template match score.")
		props["score_thresholds"] = objectProperty("Template match score thresholds.")
	default:
		return nil, false
	}
	return map[string]any{"type": "object", "properties": props, "required": []string{}, "additionalProperties": true}, true
}

func LocalAgentDockContextOutputSchema() map[string]any {
	props := localContextProperties(true)
	return strictObject(props, "skills", "dynamic_mcp", "workflow_templates", "rules")
}

func FleetAgentDockContextOutputSchema() map[string]any {
	item := contextItemSchema(false)
	warning := map[string]any{
		"type": "object", "properties": map[string]any{
			"source": stringProperty("Context section identifier."), "message": stringProperty("Safe warning message."),
		}, "required": []string{"source", "message"}, "additionalProperties": false,
	}
	rules := map[string]any{"type": "array", "items": map[string]any{"type": "string"}}
	local := strictObject(localContextProperties(false), "skills", "dynamic_mcp", "rules")
	shared := strictObject(map[string]any{
		"workflow_templates": map[string]any{"type": "array", "items": item},
		"recall": strictObject(map[string]any{
			"enabled": booleanProperty("Whether NexusDock Recall is available."),
			"items":   map[string]any{"type": "array", "items": item},
		}, "enabled", "items"),
		"rules":    rules,
		"warnings": map[string]any{"type": "array", "items": warning},
	}, "workflow_templates", "recall", "rules")
	return strictObject(map[string]any{
		"nodes": map[string]any{
			"type": "array", "description": "Enabled AgentDock nodes and their node-local context.",
			"items": strictObject(map[string]any{
				"node_id":      map[string]any{"type": "string"},
				"name":         map[string]any{"type": "string"},
				"online":       map[string]any{"type": "boolean"},
				"version":      map[string]any{"type": "string"},
				"os":           map[string]any{"type": "string"},
				"arch":         map[string]any{"type": "string"},
				"capabilities": map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
				"context":      local,
				"error":        map[string]any{"type": "string"},
			}, "node_id", "name", "online", "capabilities"),
		},
		"shared": shared,
	}, "nodes", "shared")
}

func localContextProperties(includeShared bool) map[string]any {
	skill := map[string]any{
		"type": "object", "properties": map[string]any{
			"name": stringProperty("Skill name."), "description": stringProperty("Short capability description."),
			"file": stringProperty("skill:// URI for the active SKILL.md."), "bundled": booleanProperty("Whether the Skill is bundled by AgentDock."),
		}, "required": []string{"name", "description", "file"}, "additionalProperties": false,
	}
	dynamicItem := contextItemSchema(true)
	indexItem := contextItemSchema(false)
	warning := map[string]any{
		"type": "object", "properties": map[string]any{
			"source": stringProperty("Context section identifier."), "message": stringProperty("Safe warning message."),
		}, "required": []string{"source", "message"}, "additionalProperties": false,
	}
	props := map[string]any{
		"skills":      map[string]any{"type": "array", "description": "Installed document Skill capability index.", "items": skill},
		"dynamic_mcp": map[string]any{"type": "array", "description": "Enabled dynamic MCP server capability index.", "items": dynamicItem},
		"acp": strictObject(map[string]any{
			"enabled": booleanProperty("Whether ACP is enabled."), "agent": stringProperty("Configured ACP agent name."), "description": stringProperty("Short ACP usage orientation."),
		}, "enabled", "agent", "description"),
		"rules":    map[string]any{"type": "array", "description": "Operational rules for using this AgentDock runtime.", "items": map[string]any{"type": "string"}},
		"warnings": map[string]any{"type": "array", "description": "Best-effort context sections that could not be loaded.", "items": warning},
	}
	if includeShared {
		props["workflow_templates"] = map[string]any{"type": "array", "description": "Active NexusDock Workflow template index; empty when Nexus is unavailable.", "items": indexItem}
		props["recall"] = strictObject(map[string]any{
			"enabled": booleanProperty("Whether NexusDock Recall context is configured."),
			"items":   map[string]any{"type": "array", "items": indexItem},
		}, "enabled", "items")
	}
	return props
}

func contextItemSchema(requireDescription bool) map[string]any {
	required := []string{"name"}
	if requireDescription {
		required = append(required, "description")
	}
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":        stringProperty("Capability name."),
			"description": stringProperty("Short capability description."),
		},
		"required":             required,
		"additionalProperties": false,
	}
}
