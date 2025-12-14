package ops

import (
	"agent-smith/internal/config"
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Skipping ExpandPath test: no user home dir")
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"~", home},
		{"~/", home}, // trailing slash might be handled differently depending on join, but ExpandPath joins with home
		// Wait, ExpandPath matches "~/" and joins path[2:].
		// If input is "~/", path[2:] is empty string. Join(home, "") is home.
		{"~/foo/bar", filepath.Join(home, "foo", "bar")},
		{"/abs/path", "/abs/path"},
		{"rel/path", "rel/path"},
		{"", ""},
	}

	for _, tt := range tests {
		got := ExpandPath(tt.input)
		// Normalize paths for comparison (Clean)
		if filepath.Clean(got) != filepath.Clean(tt.expected) {
			t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestApplyPersona(t *testing.T) {
	// Setup temporary directory structure
	tempDir, err := os.MkdirTemp("", "ops_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create agents (persona) directory
	agentsDir := filepath.Join(tempDir, "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a dummy agent file
	personaName := "testagent"
	agentFileName := "AGENTS.testagent.md"
	agentFilePath := filepath.Join(agentsDir, agentFileName)
	content := []byte("# Test Persona")
	if err := os.WriteFile(agentFilePath, content, 0644); err != nil {
		t.Fatal(err)
	}

	// Define Targets
	linkTarget := filepath.Join(tempDir, "LINK_ME.md")
	copyTarget := filepath.Join(tempDir, "COPY_ME.md")

	targets := []config.TargetConfig{
		{Path: linkTarget, Mode: config.TargetModeLink},
		{Path: copyTarget, Mode: config.TargetModeCopy},
	}

	// EXECUTE
	foundPath, err := ApplyPersona(personaName, []string{agentsDir, "/non/existent"}, targets)
	if err != nil {
		t.Fatalf("ApplyPersona failed: %v", err)
	}

	// VERIFY
	if foundPath != agentFilePath {
		t.Errorf("Expected found path %s, got %s", agentFilePath, foundPath)
	}

	// Verify Link
	info, err := os.Lstat(linkTarget)
	if err != nil {
		t.Errorf("Link target not created: %v", err)
	} else {
		if info.Mode()&os.ModeSymlink == 0 {
			t.Errorf("Expected symlink at %s", linkTarget)
		} else {
			dest, err := os.Readlink(linkTarget)
			if err != nil {
				t.Fatal(err)
			}
			if dest != floatPathAbs(t, agentFilePath) {
				t.Errorf("Symlink points to %s, want %s", dest, agentFilePath)
			}
		}
	}

	// Verify Copy
	info, err = os.Lstat(copyTarget)
	if err != nil {
		t.Errorf("Copy target not created: %v", err)
	} else {
		if info.Mode()&os.ModeSymlink != 0 {
			t.Errorf("Expected regular file at %s, got symlink", copyTarget)
		} else {
			readContent, err := os.ReadFile(copyTarget)
			if err != nil {
				t.Fatal(err)
			}
			if string(readContent) != string(content) {
				t.Errorf("Copy content mismatch")
			}
		}
	}
}

func floatPathAbs(t *testing.T, p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}
