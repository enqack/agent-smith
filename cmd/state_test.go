package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func TestStatePersistence(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "agents-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock HOME to point to tempDir so default paths work
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Test 1: Save state to default location (no config file used)
	viper.Reset()           // Ensure clean viper state
	viper.SetConfigFile("") // No config file

	targetFile := "/tmp/AGENTS.md"
	persona := "coder"

	if err := saveState(targetFile, persona); err != nil {
		t.Fatalf("saveState failed: %v", err)
	}

	// Verify file exists in default location
	expectedPath := filepath.Join(tempDir, ".config", "agent-smith", "status.yaml")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected status file at %s, but executed not found", expectedPath)
	}

	// Test 2: Load state
	state, err := loadState()
	if err != nil {
		t.Fatalf("loadState failed: %v", err)
	}

	if state.LastTargetFile != targetFile {
		t.Errorf("Expected LastTargetFile %s, got %s", targetFile, state.LastTargetFile)
	}
	if state.LastPersona != persona {
		t.Errorf("Expected LastPersona %s, got %s", persona, state.LastPersona)
	}

	// Test 3: Save state with custom config file
	customDir := filepath.Join(tempDir, "custom")
	if err := os.MkdirAll(customDir, 0755); err != nil {
		t.Fatal(err)
	}
	customConfig := filepath.Join(customDir, "config.yaml")
	// We don't need to create the file, just tell viper we used it
	viper.SetConfigFile(customConfig)

	newTarget := "/tmp/OTHER.md"
	newPersona := "writer"

	if err := saveState(newTarget, newPersona); err != nil {
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

	if newState.LastTargetFile != newTarget {
		t.Errorf("Expected LastTargetFile %s, got %s", newTarget, newState.LastTargetFile)
	}
}
