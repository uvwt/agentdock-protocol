package mcpcontract

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCanonicalToolContractsAreCompleteAndFresh(t *testing.T) {
	if got := len(ToolNames()); got != 8 {
		t.Fatalf("tool count = %d, want 8", got)
	}
	for _, name := range ToolNames() {
		input, ok := InputSchema(name)
		if !ok {
			t.Fatalf("missing input schema for %s", name)
		}
		if input["additionalProperties"] != false {
			t.Fatalf("%s root input must be strict: %#v", name, input["additionalProperties"])
		}
		if _, ok := AnnotationContract(name); !ok {
			t.Fatalf("missing annotations for %s", name)
		}
		if name == ToolAgentDockContext {
			continue
		}
		if _, ok := OutputSchema(name); !ok {
			t.Fatalf("missing output schema for %s", name)
		}
	}

	first, _ := InputSchema(ToolRecallSearch)
	second, _ := InputSchema(ToolRecallSearch)
	first["additionalProperties"] = true
	if reflect.DeepEqual(first, second) {
		t.Fatal("schema factories share mutable state")
	}
}

func TestInputSchemasNeverSerializeNullRequired(t *testing.T) {
	for _, name := range ToolNames() {
		schema, ok := InputSchema(name)
		if !ok {
			t.Fatalf("missing input schema for %s", name)
		}
		encoded, err := json.Marshal(schema)
		if err != nil {
			t.Fatalf("marshal %s input schema: %v", name, err)
		}
		if string(encoded) == "" || json.Valid(encoded) == false {
			t.Fatalf("%s input schema is not valid JSON: %s", name, encoded)
		}
		if required, exists := schema["required"]; exists && required == nil {
			t.Fatalf("%s input schema contains required:null: %s", name, encoded)
		}
	}

	for _, name := range []string{ToolAgentDockContext, ToolWorkspaceContext, ToolRecallMaintain} {
		schema, _ := InputSchema(name)
		if _, exists := schema["required"]; exists {
			t.Fatalf("%s should omit empty required", name)
		}
	}
}

func TestWorkspaceContextHasCanonicalDirectAndNodeProfiles(t *testing.T) {
	direct, ok := InputSchema(ToolWorkspaceContext)
	if !ok {
		t.Fatal("workspace_context direct input schema missing")
	}
	directProperties := direct["properties"].(map[string]any)
	if _, ok := directProperties["workdir"]; !ok {
		t.Fatal("workspace_context direct input is missing workdir")
	}
	if _, ok := directProperties["node_id"]; ok {
		t.Fatal("direct workspace_context must not expose node_id")
	}

	node := NodeWorkspaceContextInputSchema()
	nodeProperties := node["properties"].(map[string]any)
	if _, ok := nodeProperties["workdir"]; !ok {
		t.Fatal("node workspace_context is missing workdir")
	}
	if _, ok := nodeProperties["node_id"]; !ok {
		t.Fatal("node workspace_context is missing node_id")
	}
	required := node["required"].([]string)
	if !reflect.DeepEqual(required, []string{"node_id"}) {
		t.Fatalf("node workspace_context required = %#v", required)
	}

	output := WorkspaceContextOutputSchema()
	properties := output["properties"].(map[string]any)
	for _, name := range []string{"workdir", "workspace_root", "instructions", "workspace_skills", "warnings"} {
		if _, ok := properties[name]; !ok {
			t.Fatalf("workspace_context output missing %s", name)
		}
	}
	if output["additionalProperties"] != false {
		t.Fatalf("workspace_context output must be strict: %#v", output)
	}
	workspaceSkills := properties["workspace_skills"].(map[string]any)
	workspaceSkill := workspaceSkills["items"].(map[string]any)
	workspaceSkillProps := workspaceSkill["properties"].(map[string]any)
	for _, name := range []string{"name", "description", "file", "skill_ref", "source_type", "source_id"} {
		if _, ok := workspaceSkillProps[name]; !ok {
			t.Fatalf("workspace Skill provenance is missing %s", name)
		}
	}
	if got := workspaceSkillProps["source_type"].(map[string]any)["enum"]; !reflect.DeepEqual(got, []string{"workspace"}) {
		t.Fatalf("workspace Skill source_type enum = %#v", got)
	}
}

func TestWorkflowRootIsStrictButTemplatePayloadIsOpen(t *testing.T) {
	schema, _ := InputSchema(ToolWorkflowTemplateManage)
	props := schema["properties"].(map[string]any)
	template := props["template"].(map[string]any)
	if schema["additionalProperties"] != false || template["additionalProperties"] != true {
		t.Fatalf("unexpected workflow strictness root=%#v template=%#v", schema["additionalProperties"], template["additionalProperties"])
	}
}

func TestRecallFactsAcceptRuntimeCoercibleValues(t *testing.T) {
	schema, _ := InputSchema(ToolRecallWrite)
	facts := schema["properties"].(map[string]any)["facts"].(map[string]any)
	if facts["additionalProperties"] != true {
		t.Fatalf("facts must match runtime fmt.Sprint coercion: %#v", facts)
	}
}

func TestPrivateNoteMissingEncryptedIsStringArray(t *testing.T) {
	schema, _ := OutputSchema(ToolPrivateNoteManage)
	missing := schema["properties"].(map[string]any)["missing_encrypted"].(map[string]any)
	items := missing["items"].(map[string]any)
	if missing["type"] != "array" || items["type"] != "string" {
		t.Fatalf("missing_encrypted schema = %#v", missing)
	}
}

func TestContextHasExplicitLocalAndFleetProfiles(t *testing.T) {
	local := LocalAgentDockContextOutputSchema()
	fleet := FleetAgentDockContextOutputSchema()
	localProperties := local["properties"].(map[string]any)
	if _, ok := localProperties["skills"]; !ok {
		t.Fatal("local context is missing skills")
	}
	commonSkills, ok := localProperties["common_skills"].(map[string]any)
	if !ok {
		t.Fatal("local context is missing common_skills")
	}
	commonProperties := commonSkills["properties"].(map[string]any)
	for _, name := range []string{"root", "total", "truncated", "items"} {
		if _, ok := commonProperties[name]; !ok {
			t.Fatalf("common_skills is missing %s", name)
		}
	}
	managedSkill := localProperties["skills"].(map[string]any)["items"].(map[string]any)
	managedSkillProps := managedSkill["properties"].(map[string]any)
	for _, name := range []string{"name", "description", "file", "skill_ref", "source_type", "source_id", "content_digest"} {
		if _, ok := managedSkillProps[name]; !ok {
			t.Fatalf("managed Skill provenance is missing %s", name)
		}
	}
	if got := managedSkillProps["source_type"].(map[string]any)["enum"]; !reflect.DeepEqual(got, []string{"managed"}) {
		t.Fatalf("managed Skill source_type enum = %#v", got)
	}
	commonSkill := commonProperties["items"].(map[string]any)["items"].(map[string]any)
	commonSkillProps := commonSkill["properties"].(map[string]any)
	for _, name := range []string{"name", "description", "file", "skill_ref", "source_type", "source_id", "content_digest"} {
		if _, ok := commonSkillProps[name]; !ok {
			t.Fatalf("common Skill provenance is missing %s", name)
		}
	}
	if got := commonSkillProps["source_type"].(map[string]any)["enum"]; !reflect.DeepEqual(got, []string{"shared"}) {
		t.Fatalf("common Skill source_type enum = %#v", got)
	}
	for _, required := range local["required"].([]string) {
		if required == "common_skills" {
			t.Fatal("common_skills must remain optional for rolling compatibility with older AgentDock nodes")
		}
	}
	nodes := fleet["properties"].(map[string]any)["nodes"].(map[string]any)
	node := nodes["items"].(map[string]any)
	nodeContext := node["properties"].(map[string]any)["context"].(map[string]any)
	if _, ok := nodeContext["properties"].(map[string]any)["common_skills"]; !ok {
		t.Fatal("fleet node context is missing common_skills")
	}
	runtimeSchema, ok := localProperties["runtime"].(map[string]any)
	if !ok {
		t.Fatal("local context is missing runtime")
	}
	runtimeProperties := runtimeSchema["properties"].(map[string]any)
	for _, name := range []string{"version", "os", "arch", "agentdock_home", "agentdock_default_dir", "default_cwd", "path_model"} {
		if _, ok := runtimeProperties[name]; !ok {
			t.Fatalf("local runtime is missing %s", name)
		}
	}
	if _, ok := fleet["properties"].(map[string]any)["nodes"]; !ok {
		t.Fatal("fleet context is missing nodes")
	}
	dynamicMCP := localProperties["dynamic_mcp"].(map[string]any)
	dynamicItem := dynamicMCP["items"].(map[string]any)
	dynamicProperties := dynamicItem["properties"].(map[string]any)
	for _, name := range []string{"name", "description", "status", "tool_count", "last_error_code"} {
		if _, ok := dynamicProperties[name]; !ok {
			t.Fatalf("dynamic MCP context item is missing %s", name)
		}
	}
	if dynamicItem["additionalProperties"] != false {
		t.Fatalf("dynamic MCP context item must remain strict: %#v", dynamicItem)
	}
	if reflect.DeepEqual(local, fleet) {
		t.Fatal("local and fleet context profiles unexpectedly match")
	}
}

func TestRecallWriteBehaviorVectorsCoverSafetyBoundary(t *testing.T) {
	cases := RecallWriteBehaviorCases()
	seen := map[string]RecallWriteBehaviorCase{}
	for _, c := range cases {
		seen[c.Name] = c
		if c.DryRun && c.Expected == RecallWriteMutation {
			t.Fatalf("dry run case mutates: %#v", c)
		}
		wantsOverwrite := c.Target == "markdown" && (c.Action == "replace" || c.Action == "append" || c.Action == "patch" || c.Action == "update_fact")
		if c.OverwriteSemantic != wantsOverwrite {
			t.Fatalf("overwrite semantic mismatch for %q: got=%t want=%t", c.Name, c.OverwriteSemantic, wantsOverwrite)
		}
	}
	for _, name := range []string{
		"markdown inbox create unconfirmed mutates",
		"markdown protected create unconfirmed errors",
		"markdown replace unconfirmed previews",
		"markdown delete unconfirmed errors",
		"markdown plan previews",
		"card create confirmed dry run previews",
	} {
		if _, ok := seen[name]; !ok {
			t.Fatalf("missing behavior vector %q", name)
		}
	}
}
