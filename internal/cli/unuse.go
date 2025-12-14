package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"agent-smith/internal/ops"
	"agent-smith/internal/state"
)

// unuseCmd represents the unuse command
var unuseCmd = &cobra.Command{
	Use:   "unuse",
	Short: "Remove all configured agent targets",
	Long:  `Remove all files or links configured as targets for the agent persona.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load config to know what to remove
		// We use Cfg.Targets (populated by initConfig)
		// Or should we use state?
		// Better to use config, as that is what drives the "use" command.
		// If we use state, we might miss things if config changed.
		// If we use config, we might remove things that weren't managed by us if paths collide?
		// "unuse" usually implies reversing "use".
		// Let's iterate over Cfg.Targets.

		targets := Cfg.Targets
		if len(targets) == 0 {
			// Fallback to state if config is empty (e.g. relying on defects?)
			// But duplicate logic from status/use might be needed.
			// Let's rely on Cfg.Targets as initConfig ensures backward compatibility.
			fmt.Println("No targets configured to remove.")
			return
		}

		removedCount := 0
		errCount := 0

		for _, target := range targets {
			targetPath := ops.ExpandPath(target.Path)

			// Check if exists
			if _, err := os.Lstat(targetPath); err != nil {
				if os.IsNotExist(err) {
					// Already gone, skip
					continue
				}
				fmt.Printf("Error accessing %s: %v\n", targetPath, err)
				errCount++
				continue
			}

			// Remove
			if err := os.Remove(targetPath); err != nil {
				fmt.Printf("Error removing %s: %v\n", targetPath, err)
				errCount++
			} else {
				fmt.Printf("Removed: %s\n", targetPath)
				removedCount++
			}
		}

		// Update state to clear last persona?
		// We should probably clear the last persona in status.yaml
		// to reflect that no persona is active.
		// However, we still might want to know what targets we managed?
		// Let's clear LastPersona but keep Targets (or clear them too?)
		// If we clear targets, status will show nothing.
		// If we keep targets, status will show MISSING.
		// Usually "unuse" means "I don't want an agent active".
		// So clearing LastPersona is correct.
		// Clearning Targets logic: if we don't have targets in state, status won't verify them.
		// This seems correct for "unuse".

		// We should probably clear the last persona in status.yaml
		// to reflect that no persona is active.
		// With new schema, we set CanonicalTarget to empty.

		// In new schema, we might want to keep the AgentFile entry but just ensure targets are gone?
		// But SaveState appends/updates.
		// If we want to "unuse" everything, we probably just want to clear CanonicalTarget.
		// But existing SaveState logic loads state and modifies it.
		// If we pass empty strings, SaveState might create an empty entry.

		// Let's manually load and clear CanonicalTarget.
		st, err := state.LoadState()
		if err == nil && st != nil {
			st.CanonicalTarget = "" // Clear the active symlink pointer
			st.AgentFiles = nil     // Clear managed personas

			if err := state.WriteState(st); err != nil {
				fmt.Printf("Warning: Failed to update state file: %v\n", err)
			}
		}

		if removedCount > 0 {
			fmt.Printf("Successfully removed %d targets.\n", removedCount)
		} else if errCount == 0 {
			fmt.Println("No targets needed removal.")
		}
	},
}

func init() {
	rootCmd.AddCommand(unuseCmd)
}
