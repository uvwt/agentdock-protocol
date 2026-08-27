package mcpcontract

func stringProperty(description string) map[string]any {
	return map[string]any{"type": "string", "description": description}
}
func integerProperty(description string) map[string]any {
	return map[string]any{"type": "integer", "description": description}
}
func boundedIntegerProperty(description string, minimum, maximum int) map[string]any {
	return map[string]any{"type": "integer", "description": description, "minimum": minimum, "maximum": maximum}
}
func booleanProperty(description string) map[string]any {
	return map[string]any{"type": "boolean", "description": description}
}
func enumProperty(description string, values ...string) map[string]any {
	return map[string]any{"type": "string", "description": description, "enum": values}
}
func objectProperty(description string) map[string]any {
	return map[string]any{"type": "object", "description": description, "additionalProperties": true}
}
func arrayObjects(description string) map[string]any {
	return map[string]any{"type": "array", "description": description, "items": map[string]any{"type": "object", "additionalProperties": true}}
}
func arrayStrings(description string) map[string]any {
	return map[string]any{"type": "array", "description": description, "items": map[string]any{"type": "string"}}
}
func strictObject(properties map[string]any, required ...string) map[string]any {
	schema := map[string]any{"type": "object", "properties": properties, "additionalProperties": false}
	if len(required) > 0 {
		schema["required"] = append([]string(nil), required...)
	}
	return schema
}
func openObject(description string) map[string]any {
	return map[string]any{"type": "object", "description": description, "additionalProperties": true}
}
