package runtimecontract

import "testing"

func TestCurrentVersion(t *testing.T) {
	if CurrentVersion != 1 {
		t.Fatalf("CurrentVersion = %d, want 1", CurrentVersion)
	}
}
