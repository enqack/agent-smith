package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current agent status",
	Long:  `Show which agent persona is currently active by checking the AGENTS.md symlink.`,
	Run: func(cmd *cobra.Command, args []string) {
		targetFile := viper.GetString("target_file")

		// If the user didn't explicitly set the flag, try to load from state
		if !cmd.Flags().Changed("target-file") {
			if state, err := loadState(); err == nil && state.LastTargetFile != "" {
				targetFile = state.LastTargetFile
			}
		}

		// Check if file exists
		info, err := os.Lstat(targetFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("No agent currently selected (File %s does not exist)\n", targetFile)
				return
			}
			fmt.Printf("Error checking status: %v\n", err)
			return
		}

		if info.Mode()&os.ModeSymlink != 0 {
			linkDest, err := os.Readlink(targetFile)
			if err != nil {
				fmt.Printf("Error reading link: %v\n", err)
				return
			}
			fmt.Printf("Current Agent: %s -> %s\n", targetFile, linkDest)

			// Try to resolve the persona name
			base := filepath.Base(linkDest)
			// e.g., AGENTS.coder.md
			if len(base) > 10 && base[:7] == "AGENTS." && base[len(base)-3:] == ".md" {
				persona := base[7 : len(base)-3]
				fmt.Printf("Active Persona: %s\n", persona)
			}

		} else {
			fmt.Printf("%s exists but is not a symlink.\n", targetFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
