package e2e_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigPrecedence(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "agents-e2e-config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	agentsDir := filepath.Join(tempDir, "agents")
	os.MkdirAll(agentsDir, 0755)

	targetFile := filepath.Join(tempDir, "AGENTS.md")
	targetEnv := filepath.Join(tempDir, "ENV_TARGET.md")
	targetFlag := filepath.Join(tempDir, "FLAG_TARGET.md")

	os.WriteFile(filepath.Join(agentsDir, "AGENTS.test.md"), []byte("Test."), 0644)

	// 1. Env Var vs Config File
	// Create config file pointing to targetFile
	configDir := filepath.Join(tempDir, ".config", "agent-smith")
	os.MkdirAll(configDir, 0755)
	configFile := filepath.Join(configDir, "config.yaml")
	os.WriteFile(configFile, []byte(fmt.Sprintf("agents_dir: ['%s']\ntarget_file: '%s'", agentsDir, targetFile)), 0644)

	// Set env var. root.go uses Viper with SetEnvPrefix("AGENTS") and ReplaceDash(underscore).
	// So "target-file" -> "AGENTS_TARGET_FILE".
	os.Setenv("AGENTS_TARGET_FILE", targetEnv)
	defer os.Unsetenv("AGENTS_TARGET_FILE")

	// We MUST pass --config explicitly if not in expected location, or ensure HOME is set correct (it is in runAgentsS)
	// runAgentsS sets HOME to tempDir.
	// root.go looks in $XDG_CONFIG_HOME/agent-smith/config.yaml or ~/.config/...
	// We set XDG_CONFIG_HOME to tempDir/.config in runAgentsS.
	// So config should be loaded.

	// Run 'use test'
	out, err := runAgentsS(t, tempDir, "use", "test")
	if err != nil {
		t.Fatalf("use (env) failed: %v\nOutput: %s", err, out)
	}

	// Verify ENV_TARGET.md was used (Env > Config)
	if _, err := os.Stat(targetEnv); err != nil {
		t.Errorf("Expected env var target %s to be created, but missing. Output:\n%s", targetEnv, out)
	}
	if _, err := os.Stat(targetFile); !os.IsNotExist(err) {
		t.Errorf("Expected config target %s to NOT be created (overridden by env)", targetFile)
	}

	// 2. Flag vs Env Var
	// Run 'use test --target-file FLAG_TARGET.md'
	// Env var AGENTS_TARGET_FILE is still set to targetEnv

	// Reset targets
	os.Remove(targetEnv)

	out, err = runAgentsS(t, tempDir, "use", "test", "--target-file", targetFlag)
	if err != nil {
		t.Fatalf("use (flag) failed: %v\nOutput: %s", err, out)
	}

	// Verify FLAG target created
	if _, err := os.Stat(targetFlag); err != nil {
		t.Errorf("Expected flag target %s to be created, but missing. Output:\n%s", targetFlag, out)
	}

	// Verify ENV target ALSO created (Additive behavior as of v0.3.1)
	// Because Config/Env is the "Canonical" reference.
	if _, err := os.Stat(targetEnv); err != nil {
		t.Errorf("Expected env target %s to be created (system active), but missing", targetEnv)
	}
}
