package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestInitConfigDefaults(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "agents-root-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock HOME and XDG vars
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))
	os.Setenv("XDG_DATA_HOME", filepath.Join(tempDir, ".local", "share"))
	defer os.Unsetenv("XDG_CONFIG_HOME")
	defer os.Unsetenv("XDG_DATA_HOME")

	// Reset viper
	viper.Reset()
	cfgFile = "" // Reset configuration file flag

	// Run initConfig
	initConfig()

	// Verify target_file default
	expectedTarget := filepath.Join(tempDir, ".config", "agents", "AGENTS.md")
	if val := viper.GetString("target_file"); val != expectedTarget {
		t.Errorf("Expected target_file %s, got %s", expectedTarget, val)
	}

	// Verify agents_dir default
	// It should contain ~/.config/agent-smith/personas and /usr/share/agent-smith/personas
	agentsDirs := viper.GetStringSlice("agents_dir")
	if len(agentsDirs) < 2 {
		t.Errorf("Expected at least 2 default agents dirs, got %d", len(agentsDirs))
	}

	expectedUserAgents := filepath.Join(tempDir, ".local", "share", "agent-smith", "personas")
	foundUser := false
	for _, dir := range agentsDirs {
		if dir == expectedUserAgents {
			foundUser = true
			break
		}
	}
	if !foundUser {
		t.Errorf("Expected agents_dir to contain %s, got %v", expectedUserAgents, agentsDirs)
	}
}

func TestInitConfigCreateDirs(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "agents-root-test-dirs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock HOME and XDG vars
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, ".config"))
	defer os.Unsetenv("XDG_CONFIG_HOME")

	// Reset viper
	viper.Reset()
	cfgFile = ""

	// Run initConfig
	initConfig()

	// Verify default config directory was created: ~/.config/agent-smith
	configDir := filepath.Join(tempDir, ".config", "agent-smith")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Expected config directory %s was not created", configDir)
	}
}
