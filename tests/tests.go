package tests

import (
	"os"
	"testing"
)

// TestMain -- test entrypoint and setup
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
