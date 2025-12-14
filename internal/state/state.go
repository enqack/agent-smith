package state

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"agent-smith/internal/config"
)

type TargetState struct {
	Path string            `yaml:"path"`
	Mode config.TargetMode `yaml:"mode"`
}

type AgentFileState struct {
	Name    string        `yaml:"name"` // Determine if we still need "Name" (persona label). User said: "name (persona label)"
	Path    string        `yaml:"path"` // The actual agent_file path
	Targets []TargetState `yaml:"targets"`
}

type StatusState struct {
	CanonicalTarget string           `yaml:"canonical_target,omitempty"`
	AgentFiles      []AgentFileState `yaml:"agent_files"`
}

func getStatusFilePath() (string, error) {
	// If a config file was used, store status.yaml in the same directory
	if configFile := viper.ConfigFileUsed(); configFile != "" {
		return filepath.Join(filepath.Dir(configFile), "status.yaml"), nil
	}

	// Fallback to XDG State Home
	stateHome, err := config.GetStateHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(stateHome, "agent-smith", "status.yaml"), nil
}

func LoadState() (*StatusState, error) {
	path, err := getStatusFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state StatusState
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func SaveState(canonicalTarget, personaName, agentFile string, targets []config.TargetConfig) error {
	// Load existing state to preserve other agent files
	state, err := LoadState()
	if err != nil {
		state = &StatusState{}
	}
	if state == nil {
		state = &StatusState{}
	}

	state.CanonicalTarget = canonicalTarget

	// Convert config targets to state targets
	var stateTargets []TargetState
	for _, t := range targets {
		stateTargets = append(stateTargets, TargetState{
			Path: t.Path,
			Mode: t.Mode,
		})
	}

	// Update or Append AgentFile
	found := false
	for i, af := range state.AgentFiles {
		// Key by AgentFile Path
		if af.Path == agentFile {
			state.AgentFiles[i].Name = personaName     // Update label if it changed
			state.AgentFiles[i].Targets = stateTargets // Replace targets (Authority: "use" command)
			found = true
			break
		}
	}
	if !found && agentFile != "" {
		state.AgentFiles = append(state.AgentFiles, AgentFileState{
			Name:    personaName,
			Path:    agentFile,
			Targets: stateTargets,
		})
	}

	path, err := getStatusFilePath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(&state)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// WriteState saves the given state to the status file.
// This allows consumers to perform custom state updates (like clearing fields).
func WriteState(state *StatusState) error {
	path, err := getStatusFilePath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
