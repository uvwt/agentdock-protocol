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

	for _, name := range []string{ToolAgentDockContext, ToolRecallBootstrap, ToolRecallMaintain} {
		schema, _ := InputSchema(name)
		if _, exists := schema["required"]; exists {
			t.Fatalf("%s should omit empty required", name)
		}
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
	if _, ok := local["properties"].(map[string]any)["skills"]; !ok {
		t.Fatal("local context is missing skills")
	}
	if _, ok := fleet["properties"].(map[string]any)["nodes"]; !ok {
		t.Fatal("fleet context is missing nodes")
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
