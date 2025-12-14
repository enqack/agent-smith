package e2e_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to run agents command
func runAgentsS(t *testing.T, homeDir string, args ...string) (string, error) {
	cmd := exec.Command(testBinaryPath, args...)
	cmd.Env = append(os.Environ(), "HOME="+homeDir)
	// Ensure XDG vars are unset or pointed to sandbox to avoid leaking
	cmd.Env = append(cmd.Env, "XDG_CONFIG_HOME="+filepath.Join(homeDir, ".config"))
	cmd.Env = append(cmd.Env, "XDG_DATA_HOME="+filepath.Join(homeDir, ".local", "share"))
	cmd.Env = append(cmd.Env, "XDG_STATE_HOME="+filepath.Join(homeDir, ".local", "state"))

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestIntegrationFlow(t *testing.T) {
	// Setup isolation environment
	tempDir, err := os.MkdirTemp("", "agents-e2e-flow")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Directories
	agentsDir := filepath.Join(tempDir, "agents")
	configDir := filepath.Join(tempDir, ".config", "agent-smith")

	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
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

	// 1. List
	out, err := runAgentsS(t, tempDir, "list")
	if err != nil {
		t.Fatalf("list failed: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(out, "coder") {
		t.Errorf("list output missing 'coder':\n%s", out)
	}

	// 2. Use
	out, err = runAgentsS(t, tempDir, "use", "coder")
	if err != nil {
		t.Fatalf("use failed: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(out, "Persona switched: coder") {
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
	out, err = runAgentsS(t, tempDir, "status")
	if err != nil {
		t.Fatalf("status failed: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(out, "Persona: coder [ACTIVE]") {
		t.Errorf("status output unexpected:\n%s", out)
	}
}

func TestMultiTargetFeatures(t *testing.T) {
	// Setup isolation environment
	tempDir, err := os.MkdirTemp("", "agents-e2e-multi")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Directories
	agentsDir := filepath.Join(tempDir, "agents")
	configDir := filepath.Join(tempDir, ".config", "agent-smith")
	targetsDir := filepath.Join(tempDir, "targets")

	os.MkdirAll(agentsDir, 0755)
	os.MkdirAll(configDir, 0755)
	os.MkdirAll(targetsDir, 0755)

	// Create dummy personas
	os.WriteFile(filepath.Join(agentsDir, "AGENTS.coder.md"), []byte("Code."), 0644)
	os.WriteFile(filepath.Join(agentsDir, "AGENTS.writer.md"), []byte("Write."), 0644)

	targetLink := filepath.Join(targetsDir, "LINK.md")
	targetCopy := filepath.Join(targetsDir, "COPY.md")

	configFile := filepath.Join(configDir, "config.yaml")

	// Create config with multiple targets
	configContent := fmt.Sprintf(`
agents_dir: ["%s"]
target_file: "%s"
targets:
  - path: "%s"
    mode: "link"
  - path: "%s"
    mode: "copy"
`, agentsDir, targetLink, targetLink, targetCopy)

	os.WriteFile(configFile, []byte(configContent), 0644)

	// 1. Use 'coder'
	out, err := runAgentsS(t, tempDir, "use", "coder")
	if err != nil {
		t.Fatalf("use failed: %v\nOutput: %s", err, out)
	}

	// Verify Link
	link, err := os.Readlink(targetLink)
	if err != nil {
		t.Fatalf("Failed to check link: %v", err)
	}
	if filepath.Base(link) != "AGENTS.coder.md" {
		t.Errorf("Expected link to AGENTS.coder.md, got %s", link)
	}

	// Verify Copy
	content, err := os.ReadFile(targetCopy)
	if err != nil {
		t.Fatalf("Failed to read copy: %v", err)
	}
	if string(content) != "Code." {
		t.Errorf("Expected copy content 'Code.', got '%s'", string(content))
	}

	// 2. Status - Should be OK
	out, err = runAgentsS(t, tempDir, "status")
	if err != nil {
		t.Fatalf("status failed: %v\nOutput: %s", err, out)
	}
	if !strings.Contains(out, "Persona: coder [ACTIVE]") {
		t.Errorf("status missing active persona:\n%s", out)
	}
	if !strings.Contains(out, "[OK]") {
		t.Errorf("status missing OK:\n%s", out)
	}

	// 3. Introduces Drift in Link
	os.Remove(targetLink)
	os.Symlink(filepath.Join(agentsDir, "AGENTS.writer.md"), targetLink)

	out, err = runAgentsS(t, tempDir, "status")
	// With new logic, changing canonical target IS switching persona.
	// So we expect Active Persona to be writer based on canonical link.
	if !strings.Contains(out, "Persona: writer [ACTIVE]") {
		t.Errorf("status expected Active Persona: writer, got:\n%s", out)
	}

	// 4. Reconcile
	// Note: Reconcile now respects the ACTIVE persona defined by the Canonical Target.
	// Since we manually switched to "writer", Reconcile accepts "writer" as the source of truth
	// and should update other targets to match "writer".

	out, err = runAgentsS(t, tempDir, "reconcile")
	if err != nil {
		t.Fatalf("reconcile failed: %v\nOutput: %s", err, out)
	}

	// Verify Link remains Writer
	link, _ = os.Readlink(targetLink)
	if filepath.Base(link) != "AGENTS.writer.md" {
		t.Errorf("Reconcile should respect canonical drift. Expected writer, got %s", link)
	}

	// Verify Copy updated to Writer (Propagated Truth)
	content, err = os.ReadFile(targetCopy)
	if err != nil {
		t.Fatalf("Failed to read copy: %v", err)
	}
	if string(content) != "Write." {
		t.Errorf("Expected copy content to update to 'Write.', got '%s'", string(content))
	}
}
