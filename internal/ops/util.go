package ops

import (
	"os"
	"path/filepath"
)

// ExpandPath expands ~ to the home directory
func ExpandPath(path string) string {
	if path == "~" {
		home, _ := os.UserHomeDir()
		return home
	}
	if len(path) > 1 && path[:2] == "~/" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
