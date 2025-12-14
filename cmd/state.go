package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type StatusState struct {
	LastTargetFile string `yaml:"last_target_file"`
	LastPersona    string `yaml:"last_persona"`
}

func getStatusFilePath() (string, error) {
	// If a config file was used, store status.yaml in the same directory
	if configFile := viper.ConfigFileUsed(); configFile != "" {
		return filepath.Join(filepath.Dir(configFile), "status.yaml"), nil
	}

	// Fallback to XDG State Home
	stateHome, err := getStateHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(stateHome, "agent-smith", "status.yaml"), nil
}

func loadState() (*StatusState, error) {
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

func saveState(targetFile, persona string) error {
	path, err := getStatusFilePath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	state := StatusState{
		LastTargetFile: targetFile,
		LastPersona:    persona,
	}

	data, err := yaml.Marshal(&state)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
