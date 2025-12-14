package config

import (
	"os"
	"path/filepath"
)

// GetConfigHome returns the XDG_CONFIG_HOME or default ~/.config
func GetConfigHome() (string, error) {
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return xdgConfigHome, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config"), nil
}

// GetDataHome returns the XDG_DATA_HOME or default ~/.local/share
func GetDataHome() (string, error) {
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return xdgDataHome, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share"), nil
}

// GetStateHome returns the XDG_STATE_HOME or default ~/.local/state
func GetStateHome() (string, error) {
	if xdgStateHome := os.Getenv("XDG_STATE_HOME"); xdgStateHome != "" {
		return xdgStateHome, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "state"), nil
}
