package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [persona]",
	Short: "Switch to a specific persona",
	Long: `Switch the current AGENTS.md symlink to point to the specified persona.
Example: agents use coder`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		persona := args[0]
		agentsDirs := viper.GetStringSlice("agents_dir")
		if len(agentsDirs) == 0 {
			if s := viper.GetString("agents_dir"); s != "" {
				agentsDirs = []string{s}
			}
		}
		targetFile := viper.GetString("target_file")

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
			return
		}

		// Ensure target directory exists
		targetDir := filepath.Dir(targetFile)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			fmt.Printf("Error creating target directory %s: %v\n", targetDir, err)
			return
		}

		// 2. Remove existing target file/symlink
		if _, err := os.Lstat(targetFile); err == nil {
			err = os.Remove(targetFile)
			if err != nil {
				fmt.Printf("Error removing existing %s: %v\n", targetFile, err)
				return
			}
		}

		// 3. Create new symlink
		// Note: We use absolute path for the source to avoid relative path confusion
		absAgentPath, err := filepath.Abs(agentPath)
		if err != nil {
			// Fallback to configured path if Abs fails for some reason
			absAgentPath = agentPath
		}

		err = os.Symlink(absAgentPath, targetFile)
		if err != nil {
			fmt.Printf("Error creating symlink: %v\n", err)
			return
		}

		// Save state for 'status' command
		if err := saveState(targetFile, persona); err != nil {
			// Don't fail the operation, but warn user
			fmt.Printf("Warning: Failed to save status state: %v\n", err)
		}

		fmt.Println("The mind was never changed; only where it points.")
		fmt.Printf("Persona switched: %s\n", persona)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
