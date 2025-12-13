package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestListCommand(t *testing.T) {
	// Setup
	t.Cleanup(func() { viper.Reset() })

	tmpDir, err := os.MkdirTemp("", "agents-test-list")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create dummy agents
	os.WriteFile(filepath.Join(tmpDir, "AGENTS.coder.md"), []byte("coder agent"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "AGENTS.writer.md"), []byte("writer agent"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "AGENTS.md"), []byte("default agent"), 0644) // Should be ignored
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("readme"), 0644)       // Should be ignored

	// Capture output
	buf := new(bytes.Buffer)
	listCmd.SetOut(buf)
	listCmd.SetErr(buf)

	// Set flags
	viper.Set("agents_dir", []string{tmpDir})

	// Run command
	listCmd.Run(listCmd, []string{})

	output := buf.String()
	assert.Contains(t, output, "coder")
	assert.Contains(t, output, "writer")
	assert.NotContains(t, output, "README")
	// "AGENTS.md" (the base one) is explicitly excluded in list.go
	// The list command extracts persona names, so if AGENTS.md was processed incorrectly as a persona,
	// it would likely show up as just "AGENTS" or similar, but definitely not with .md extension in the listing part.
	// However, let's just ensure we don't see it listed as an agent.
	assert.NotContains(t, output, "- AGENTS (")
}

func TestUseCommand(t *testing.T) {
	// Setup
	t.Cleanup(func() { viper.Reset() })

	tmpDir, err := os.MkdirTemp("", "agents-test-use")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	agentsDir := filepath.Join(tmpDir, "agents")
	err = os.Mkdir(agentsDir, 0755)
	assert.NoError(t, err)

	configDir := filepath.Join(tmpDir, "config")
	err = os.Mkdir(configDir, 0755)
	assert.NoError(t, err)

	targetFile := filepath.Join(configDir, "AGENTS.md")

	// Create agent
	agentPath := filepath.Join(agentsDir, "AGENTS.coder.md")
	os.WriteFile(agentPath, []byte("coder agent content"), 0644)

	// Capture output
	buf := new(bytes.Buffer)
	useCmd.SetOut(buf)
	useCmd.SetErr(buf)

	// Set viper config
	viper.Set("agents_dir", []string{agentsDir})
	viper.Set("target_file", targetFile)

	// Run command
	useCmd.Run(useCmd, []string{"coder"})

	// Verify symlink
	info, err := os.Lstat(targetFile)
	assert.NoError(t, err)
	assert.True(t, info.Mode()&os.ModeSymlink != 0, "Target should be a symlink")

	dest, err := os.Readlink(targetFile)
	assert.NoError(t, err)

	// Compare absolute paths
	absAgentPath, _ := filepath.Abs(agentPath)
	assert.Equal(t, absAgentPath, dest)

	// Verify content
	content, err := os.ReadFile(targetFile)
	assert.NoError(t, err)
	assert.Equal(t, "coder agent content", string(content))
}
