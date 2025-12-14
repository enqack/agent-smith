package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestXDGFallbacks(t *testing.T) {
	// Mock HOME
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// Unset XDG vars to test fallbacks
	os.Unsetenv("XDG_CONFIG_HOME")
	os.Unsetenv("XDG_DATA_HOME")
	os.Unsetenv("XDG_STATE_HOME")

	// Config Home -> ~/.config
	got, err := GetConfigHome()
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(tempDir, ".config")
	if got != expected {
		t.Errorf("GetConfigHome() = %s, want %s", got, expected)
	}

	// Data Home -> ~/.local/share
	got, err = GetDataHome()
	if err != nil {
		t.Fatal(err)
	}
	expected = filepath.Join(tempDir, ".local", "share")
	if got != expected {
		t.Errorf("GetDataHome() = %s, want %s", got, expected)
	}

	// State Home -> ~/.local/state
	got, err = GetStateHome()
	if err != nil {
		t.Fatal(err)
	}
	expected = filepath.Join(tempDir, ".local", "state")
	if got != expected {
		t.Errorf("GetStateHome() = %s, want %s", got, expected)
	}
}

func TestXDGOverrides(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_test_xdg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Set XDG vars
	xdgConfig := filepath.Join(tempDir, "custom_config")
	xdgData := filepath.Join(tempDir, "custom_data")
	xdgState := filepath.Join(tempDir, "custom_state")

	os.Setenv("XDG_CONFIG_HOME", xdgConfig)
	os.Setenv("XDG_DATA_HOME", xdgData)
	os.Setenv("XDG_STATE_HOME", xdgState)

	defer func() {
		os.Unsetenv("XDG_CONFIG_HOME")
		os.Unsetenv("XDG_DATA_HOME")
		os.Unsetenv("XDG_STATE_HOME")
	}()

	// Config Home
	got, err := GetConfigHome()
	if err != nil {
		t.Fatal(err)
	}
	if got != xdgConfig {
		t.Errorf("GetConfigHome() = %s, want %s", got, xdgConfig)
	}

	// Data Home
	got, err = GetDataHome()
	if err != nil {
		t.Fatal(err)
	}
	if got != xdgData {
		t.Errorf("GetDataHome() = %s, want %s", got, xdgData)
	}

	// State Home
	got, err = GetStateHome()
	if err != nil {
		t.Fatal(err)
	}
	if got != xdgState {
		t.Errorf("GetStateHome() = %s, want %s", got, xdgState)
	}
}
