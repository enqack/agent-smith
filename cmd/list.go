package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available agent personas",
	Long: `List all available agent personas found in the configured agents directory.
Personas are defined in files named AGENTS.<persona>.md`,
	Run: func(cmd *cobra.Command, args []string) {
		agentsDirs := viper.GetStringSlice("agents_dir")
		if len(agentsDirs) == 0 {
			// Fallback to string if slice is empty (e.g. env var set as string)
			if s := viper.GetString("agents_dir"); s != "" {
				agentsDirs = []string{s}
			}
		}

		fmt.Println("Available Agents:")
		foundAny := false

		seen := make(map[string]bool)

		for _, agentsDir := range agentsDirs {
			// Ensure agents directory exists (only if it looks like a default/user one we should create?
			// Actually previous logic created it. Let's create it if it doesn't exist to be friendly,
			// or just skip if missing. The prompt asked for creation: "The program should create any required directories"
			// But I successfully updated that logic before revert... wait.
			// The user REVERTED the create logic because of NixOS read-only.
			// So I should NOT create directories here. Just skip if missing.

			files, err := os.ReadDir(agentsDir)
			if err != nil {
				if !os.IsNotExist(err) {
					// Only print error if it's not permission denied/not found?
					// fmt.Printf("Error reading %s: %v\n", agentsDir, err)
				}
				continue
			}

			for _, file := range files {
				if file.IsDir() {
					continue
				}
				name := file.Name()
				if strings.HasPrefix(name, "AGENTS.") && strings.HasSuffix(name, ".md") && name != "AGENTS.md" {
					// Extract persona name: AGENTS.coder.md -> coder
					persona := strings.TrimSuffix(strings.TrimPrefix(name, "AGENTS."), ".md")
					if !seen[persona] {
						fmt.Printf("  - %s (%s)\n", persona, agentsDir)
						seen[persona] = true
						foundAny = true
					}
				}
			}
		}

		if !foundAny {
			fmt.Println("  (No agents found)")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
