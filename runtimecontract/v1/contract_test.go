package runtimecontractv1

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestOpenAPIAdvertisesV1(t *testing.T) {
	if !strings.Contains(OpenAPI(), "title: AgentDock Runtime API Contract v1") || !strings.Contains(OpenAPI(), "  version: \"1\"") {
		t.Fatal("OpenAPI metadata does not advertise Runtime contract v1")
	}
}

func TestGeneratedStatusAndSkillModelsMatchWireNames(t *testing.T) {
	var status StatusResponse
	if err := json.Unmarshal([]byte(`{"ok":true,"source":"agentdock-api","service":"agentdock","version":"0.9.1","runtime_contract_version":1}`), &status); err != nil {
		t.Fatal(err)
	}
	if status.RuntimeContractVersion != 1 || status.Version != "0.9.1" {
		t.Fatalf("status = %#v", status)
	}
	var list SkillListResponse
	if err := json.Unmarshal([]byte(`{"action":"list","count":1,"source":"agentdock-api","skills":[{"skill":"demo","name":"demo","description":"demo","skill_ref":"skill://managed/demo","source_type":"managed","content_digest":"sha256:demo","file_count":1}]}`), &list); err != nil {
		t.Fatal(err)
	}
	if len(list.Skills) != 1 || list.Skills[0].SkillRef != "skill://managed/demo" {
		t.Fatalf("skills = %#v", list.Skills)
	}
}
