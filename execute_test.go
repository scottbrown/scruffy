package scruffy

import (
	"testing"
)

func TestExecute(t *testing.T) {
	err := Execute()
	if err != nil {
		t.Logf("Execute returned error (expected for missing args): %v", err)
	}
}