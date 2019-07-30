// !build

package tests

import (
	"os"
	"testing"
)

//go:generate gogen-enum -input ./enums.yaml -package tests -output ./enums.go
//go:generate gofmt -w enums.go
//go:generate golangci-lint run enums.go

// TestMain -- test entrypoint and setup
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
