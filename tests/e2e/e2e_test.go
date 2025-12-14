package e2e_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Global path to the built binary for all tests in this package
var testBinaryPath string

func TestMain(m *testing.M) {
	// Build the binary once
	tempDir, err := os.MkdirTemp("", "agents-build")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	testBinaryPath = filepath.Join(tempDir, "agents")

	// Determine project root (../../ from tests/e2e)
	wd, _ := os.Getwd()
	projectRoot := filepath.Join(wd, "../..")

	// Build with CGO_ENABLED=0
	cmd := exec.Command("go", "build", "-o", testBinaryPath, "./cmd/agents")
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build binary: %v\nOutput: %s\n", err, out)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}
