package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// invokeCLI runs the agents binary with the given arguments.
// It assumes the binary has been built and location is passed or we build it.
// For simplicity, we will build a temporary binary for the test run.
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
	// We need to build from the project root.
	// Since this test is in ./cmd, the root is ..
	projectRoot, _ := filepath.Abs("..")

	// Build with CGO_ENABLED=0
	cmd := exec.Command("go", "build", "-o", testBinaryPath, ".")
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

func TestIntegrationFlow(t *testing.T) {
	// Setup isolation environment
	tempDir, err := os.MkdirTemp("", "agents-integration")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Directories
	agentsDir := filepath.Join(tempDir, "agents")
	configDir := filepath.Join(tempDir, "config")
	err = os.MkdirAll(agentsDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create a dummy agent
	err = os.WriteFile(filepath.Join(agentsDir, "AGENTS.coder.md"), []byte("Start coding."), 0644)
	if err != nil {
		t.Fatal(err)
	}

	targetFile := filepath.Join(tempDir, "AGENTS.md")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create a config file
	configContent := fmt.Sprintf(`
agents_dir: ["%s"]
target_file: "%s"
`, agentsDir, targetFile)
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Helper to run command with config
	runAgents := func(args ...string) (string, error) {
		cmdArgs := append([]string{"--config", configFile}, args...)
		cmd := exec.Command(testBinaryPath, cmdArgs...)
		// Set HOME to tempDir so state is saved there if fallback happens
		cmd.Env = append(os.Environ(), "HOME="+tempDir)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}

	// 1. List
	out, err := runAgents("list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(out, "coder") {
		t.Errorf("list output missing 'coder':\n%s", out)
	}

	// 2. Use
	out, err = runAgents("use", "coder")
	if err != nil {
		t.Fatalf("use failed: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(out, "Switched to agent: coder") {
		t.Errorf("use output unexpected:\n%s", out)
	}

	// Verify symlink
	link, err := os.Readlink(targetFile)
	if err != nil {
		t.Fatalf("Failed to read link %s: %v", targetFile, err)
	}
	if filepath.Base(link) != "AGENTS.coder.md" {
		t.Errorf("Symlink points to %s, expected AGENTS.coder.md", link)
	}

	// 3. Status
	out, err = runAgents("status")
	if err != nil {
		t.Fatalf("status failed: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(out, "Active Persona: coder") {
		t.Errorf("status output unexpected:\n%s", out)
	}
}
