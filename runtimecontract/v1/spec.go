package runtimecontractv1

import _ "embed"

// openAPISpec is the canonical v1 Runtime wire contract used for generation and validation.
//
//go:embed openapi.yaml
var openAPISpec string

// OpenAPI returns the immutable OpenAPI source for Runtime contract v1.
func OpenAPI() string { return openAPISpec }
