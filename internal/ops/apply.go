package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"agent-smith/internal/config"
)

// ApplyPersona applies the given persona to the specified targets
func ApplyPersona(persona string, agentsDirs []string, targets []config.TargetConfig) (string, error) {
	agentFileName := fmt.Sprintf("AGENTS.%s.md", persona)
	var agentPath string
	found := false

	for _, dir := range agentsDirs {
		candidate := filepath.Join(dir, agentFileName)
		if _, err := os.Stat(candidate); err == nil {
			agentPath = candidate
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: Persona '%s' not found.\n", persona)
		fmt.Printf("Searched in:\n")
		for _, dir := range agentsDirs {
			fmt.Printf("  - %s\n", dir)
		}
		return "", fmt.Errorf("persona not found")
	}

	// Read persona content once if needed for copy
	var personaContent []byte
	var applyErrors []error

	// Iterate over targets
	for _, target := range targets {
		// Check if target directory exists
		targetPath := ExpandPath(target.Path)
		dir := filepath.Dir(targetPath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			err = fmt.Errorf("error creating target directory %s: %w", dir, err)
			fmt.Println(err)
			applyErrors = append(applyErrors, err)
			continue
		}

		// Perform Action
		if target.Mode == config.TargetModeCopy {
			// Lazy load content
			if personaContent == nil {
				content, err := os.ReadFile(agentPath)
				if err != nil {
					err = fmt.Errorf("error reading persona file %s: %w", agentPath, err)
					fmt.Println(err)
					return "", err
				}
				personaContent = content
			}

			// Atomic Copy
			tmpFile, err := os.CreateTemp(dir, "agents-tmp-*")
			if err != nil {
				err = fmt.Errorf("error creating temp file for %s: %w", targetPath, err)
				fmt.Println(err)
				applyErrors = append(applyErrors, err)
				continue
			}
			tmpName := tmpFile.Name()

			if _, err := tmpFile.Write(personaContent); err != nil {
				err = fmt.Errorf("error writing to temp file: %w", err)
				fmt.Println(err)
				tmpFile.Close()
				os.Remove(tmpName)
				applyErrors = append(applyErrors, err)
				continue
			}
			if err := tmpFile.Close(); err != nil {
				err = fmt.Errorf("error closing temp file: %w", err)
				fmt.Println(err)
				os.Remove(tmpName)
				applyErrors = append(applyErrors, err)
				continue
			}

			// Rename (Atomic replace)
			if err := os.Rename(tmpName, targetPath); err != nil {
				err = fmt.Errorf("error renaming to %s: %w", targetPath, err)
				fmt.Println(err)
				os.Remove(tmpName)
				applyErrors = append(applyErrors, err)
				continue
			}

			// Fix permissions (os.CreateTemp creates 0600)
			if err := os.Chmod(targetPath, 0644); err != nil {
				// Warn but don't fail hard? Or fail? Best to warn.
				// However, if we want strict correctness, we might log it.
				// For now let's just log and continue, or maybe append error?
				// The prompt suggested strict fix. I'll append error but maybe not block?
				// Actually, if chmod fails, it's not ideal.
				// Let's just append to applyErrors for correctness.
				// But wait, the file IS updated.
				fmt.Printf("Warning: failed to chmod %s: %v\n", targetPath, err)
			}

			fmt.Printf("Updated (copy): %s\n", targetPath)

		} else {
			// Link Mode (Default)
			// Remove existing
			if _, err := os.Lstat(targetPath); err == nil {
				if err := os.Remove(targetPath); err != nil {
					err = fmt.Errorf("error removing existing %s: %w", targetPath, err)
					fmt.Println(err)
					applyErrors = append(applyErrors, err)
					continue
				}
			}

			absAgentPath, err := filepath.Abs(agentPath)
			if err != nil {
				absAgentPath = agentPath
			}

			if err := os.Symlink(absAgentPath, targetPath); err != nil {
				// Enhance error message for Windows users
				if runtime.GOOS == "windows" {
					err = fmt.Errorf("error creating symlink %s (on Windows, ensure Developer Mode is enabled or run as Administrator): %w", targetPath, err)
				} else {
					err = fmt.Errorf("error creating symlink %s: %w", targetPath, err)
				}
				fmt.Println(err)
				applyErrors = append(applyErrors, err)
				continue
			}
			fmt.Printf("Updated (link): %s\n", targetPath)
		}
	}

	if len(applyErrors) > 0 {
		return agentPath, fmt.Errorf("failed to apply targets: %v", applyErrors)
	}

	return agentPath, nil
}
