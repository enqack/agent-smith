package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"agent-smith/internal/config"
)

func TestStatePersistence(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "agents-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock HOME and XDG vars
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	os.Setenv("XDG_STATE_HOME", filepath.Join(tempDir, ".local", "state"))
	defer os.Unsetenv("XDG_STATE_HOME")

	// Test 1: Save state to default location (no config file used)
	viper.Reset()           // Ensure clean viper state
	viper.SetConfigFile("") // No config file

	targetPath := "/tmp/AGENTS.md"
	persona := "coder"
	targets := []config.TargetConfig{{Path: targetPath, Mode: config.TargetModeLink}}

	canonicalTarget := targetPath
	agentFile := "/tmp/agents/AGENTS.coder.md"

	if err := SaveState(canonicalTarget, persona, agentFile, targets); err != nil {
		t.Fatalf("saveState failed: %v", err)
	}

	// Verify file exists in default location
	expectedPath := filepath.Join(tempDir, ".local", "state", "agent-smith", "status.yaml")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected status file at %s, but executed not found", expectedPath)
	}

	// Test 2: Load state
	state, err := LoadState()
	if err != nil {
		t.Fatalf("loadState failed: %v", err)
	}

	// Verify CanonicalTarget
	if state.CanonicalTarget != canonicalTarget {
		t.Errorf("Expected CanonicalTarget %s, got %s", canonicalTarget, state.CanonicalTarget)
	}

	// Verify AgentFiles
	if len(state.AgentFiles) != 1 {
		t.Fatalf("Expected 1 agent file, got %d", len(state.AgentFiles))
	}
	af := state.AgentFiles[0]
	if af.Name != persona {
		t.Errorf("Expected AgentFile Name %s, got %s", persona, af.Name)
	}
	if af.Path != agentFile {
		t.Errorf("Expected AgentFile Path %s, got %s", agentFile, af.Path)
	}
	if len(af.Targets) != 1 {
		t.Fatalf("Expected 1 target, got %d", len(af.Targets))
	}
	if af.Targets[0].Path != targetPath {
		t.Errorf("Expected Target Path %s, got %s", targetPath, af.Targets[0].Path)
	}
	// TargetState.Persona removed, so no check needed

	// Test 3: Save state with custom config file
	customDir := filepath.Join(tempDir, "custom")
	if err := os.MkdirAll(customDir, 0755); err != nil {
		t.Fatal(err)
	}
	customConfig := filepath.Join(customDir, "config.yaml")
	// We don't need to create the file, just tell viper we used it
	viper.SetConfigFile(customConfig)

	newTargetPath := "/tmp/OTHER.md"
	newPersona := "writer"
	newAgentFile := "/tmp/agents/AGENTS.writer.md"
	newTargets := []config.TargetConfig{{Path: newTargetPath, Mode: config.TargetModeCopy}}

	if err := SaveState(newTargetPath, newPersona, newAgentFile, newTargets); err != nil {
		t.Fatalf("saveState (custom config) failed: %v", err)
	}

	// Verify file exists in custom location
	expectedCustomPath := filepath.Join(customDir, "status.yaml")
	if _, err := os.Stat(expectedCustomPath); os.IsNotExist(err) {
		t.Errorf("Expected status file at %s, but not found", expectedCustomPath)
	}

	// Verify content
	data, err := os.ReadFile(expectedCustomPath)
	if err != nil {
		t.Fatal(err)
	}
	var newState StatusState
	if err := yaml.Unmarshal(data, &newState); err != nil {
		t.Fatal(err)
	}

	if len(newState.AgentFiles) != 1 { // Should start fresh or append?
		// saveState loads existing. For this new dir/config, it's fresh.
		t.Fatalf("Expected 1 agent file, got %d", len(newState.AgentFiles))
	}
	af2 := newState.AgentFiles[0]
	if af2.Name != newPersona {
		t.Errorf("Expected Name %s, got %s", newPersona, af2.Name)
	}
	if af2.Targets[0].Path != newTargetPath {
		t.Errorf("Expected Target Path %s, got %s", newTargetPath, af2.Targets[0].Path)
	}
}

func TestSaveStateReplacesTargets(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "agents-replace-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	os.Setenv("XDG_STATE_HOME", filepath.Join(tempDir, "state"))
	defer os.Unsetenv("XDG_STATE_HOME")

	viper.Reset()
	viper.SetConfigFile("")

	agentFile := "/tmp/agents/AGENTS.replace.md"
	persona := "replace"

	// 1. Save first target
	target1 := "/tmp/target1.md"
	if err := SaveState(target1, persona, agentFile, []config.TargetConfig{{Path: target1, Mode: config.TargetModeLink}}); err != nil {
		t.Fatalf("First save failed: %v", err)
	}

	// 2. Save second target (effectively a new 'use' command with different target)
	target2 := "/tmp/target2.md"
	if err := SaveState(target2, persona, agentFile, []config.TargetConfig{{Path: target2, Mode: config.TargetModeCopy}}); err != nil {
		t.Fatalf("Second save failed: %v", err)
	}

	// 3. Load and verify ONLY second exists
	st, err := LoadState()
	if err != nil {
		t.Fatal(err)
	}

	if len(st.AgentFiles) != 1 {
		t.Fatalf("Expected 1 agent file, got %d", len(st.AgentFiles))
	}
	af := st.AgentFiles[0]
	// Should replaced, so length 1
	if len(af.Targets) != 1 {
		t.Fatalf("Expected 1 target (replaced), got %d: %v", len(af.Targets), af.Targets)
	}

	// Verify target is target2
	if af.Targets[0].Path != target2 {
		t.Errorf("Expected target %s, got %s", target2, af.Targets[0].Path)
	}
}
