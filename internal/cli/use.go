package cli

import (
	"agent-smith/internal/config"
	"agent-smith/internal/ops"
	"agent-smith/internal/state"
	"fmt"
	"os"

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

		// Canonical System Path (from Config/Env/Default) - defines "Active" status
		canonicalTarget := viper.GetString("target_file")

		// Operational Target Path - where we write the link/copy
		// Default to canonical, but override with flag if set
		targetPath := canonicalTarget
		if cmd.Flags().Changed("target-file") {
			targetPath, _ = cmd.Flags().GetString("target-file")
		}

		// Prepare Targets
		// We start with Cfg.Targets (GLOBAL config targets)
		// If a CLI flag target is specified, we ADD it for this run.

		targetsToApply := Cfg.Targets

		if cmd.Flags().Changed("target-file") {
			dynamicTarget := config.TargetConfig{
				Path: targetPath,
				Mode: config.TargetModeLink, // Default to link for CLI flag
			}
			targetsToApply = append(targetsToApply, dynamicTarget)
		}

		// Apply Logic ONCE
		var err error
		var agentPath string

		agentPath, err = ops.ApplyPersona(persona, agentsDirs, targetsToApply)
		if err != nil {
			// ApplyPersona prints specific errors
			os.Exit(1)
		}

		// Update Cfg.Targets in memory? No need, we used a local slice.
		// But for SaveState, we want to reflect what we just did?
		// Actually, if we added a dynamic target, should `status` track it?
		// If so, we should pass `targetsToApply` to SaveState.
		// If the user said "use --target-file foo", they expect foo to be tracked?
		// Probably yes.

		// Save state for 'status' command
		// We *always* pass the Canonical Target to SaveState, ensuring status tracks the System Active.
		if err := state.SaveState(canonicalTarget, persona, agentPath, targetsToApply); err != nil {
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
