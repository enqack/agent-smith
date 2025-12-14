package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"agent-smith/internal/ops"
	"agent-smith/internal/state"
)

// dropCmd represents the drop command
var dropCmd = &cobra.Command{
	Use:   "drop [persona]",
	Short: "Drop a target or an entire persona from tracking",
	Long: `Drop a specific target for a persona, or the entire persona if no target is specified.
This removes the target file/symlink and updates the state.

Example:
  agents drop johnny --target-file ~/test/AGENTS.md
  agents drop johnny (removes all targets for johnny)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		personaName := args[0]
		targetFile, _ := cmd.Flags().GetString("target-file")

		// Load state
		st, err := state.LoadState()
		if err != nil {
			fmt.Printf("Error loading state: %v\n", err)
			os.Exit(1)
		}
		if st == nil {
			fmt.Println("No state found.")
			return
		}

		foundPersona := false
		personaIndex := -1
		for i, af := range st.AgentFiles {
			if af.Name == personaName {
				foundPersona = true
				personaIndex = i
				break
			}
		}

		if !foundPersona {
			fmt.Printf("Persona '%s' not found in state.\n", personaName)
			return
		}

		targetsToRemove := []string{}

		if targetFile != "" {
			// Remove specific target
			targetPath := ops.ExpandPath(targetFile)

			// Find target in persona logic
			newTargets := []state.TargetState{}
			foundTarget := false

			for _, t := range st.AgentFiles[personaIndex].Targets {
				tExpanded := ops.ExpandPath(t.Path)
				if tExpanded == targetPath || t.Path == targetFile {
					foundTarget = true
					targetsToRemove = append(targetsToRemove, t.Path)
				} else {
					newTargets = append(newTargets, t)
				}
			}

			if !foundTarget {
				fmt.Printf("Target '%s' not found for persona '%s'.\n", targetFile, personaName)
				return
			}

			st.AgentFiles[personaIndex].Targets = newTargets
			fmt.Printf("Dropping target '%s' from persona '%s'...\n", targetFile, personaName)

		} else {
			// Remove ALL targets for this persona
			for _, t := range st.AgentFiles[personaIndex].Targets {
				targetsToRemove = append(targetsToRemove, t.Path)
			}

			// Remove the persona entry itself?
			// Yes, 'drop johnny' implies forgetting johnny.
			// Remove from slice
			newAgentFiles := []state.AgentFileState{}
			for i, af := range st.AgentFiles {
				if i != personaIndex {
					newAgentFiles = append(newAgentFiles, af)
				}
			}
			st.AgentFiles = newAgentFiles
			fmt.Printf("Dropping persona '%s' and all its targets...\n", personaName)
		}

		// Perform physical removal
		for _, tPath := range targetsToRemove {
			exp := ops.ExpandPath(tPath)

			// Safety Check: Is this target used by another persona?
			isUsed := false
			usedBy := ""
			for _, af := range st.AgentFiles {
				if af.Name == personaName {
					continue // clear, we are removing from this one
				}
				for _, t := range af.Targets {
					if ops.ExpandPath(t.Path) == exp {
						isUsed = true
						usedBy = af.Name
						break
					}
				}
				if isUsed {
					break
				}
			}

			if isUsed {
				fmt.Printf("State updated, but file retained: %s (Also used by '%s')\n", exp, usedBy)
			} else {
				// Check if directory
				fi, err := os.Stat(exp)
				if err == nil && fi.IsDir() {
					fmt.Printf("Warning: target '%s' is a directory. Refusing to remove.\n", exp)
					continue
				}

				if err := os.Remove(exp); err != nil {
					if !os.IsNotExist(err) {
						fmt.Printf("Warning: Failed to remove file %s: %v\n", exp, err)
					} else {
						fmt.Printf("File %s already gone.\n", exp)
					}
				} else {
					fmt.Printf("Removed: %s\n", exp)
				}
			}
		}

		// Save State
		if err := state.WriteState(st); err != nil {
			fmt.Printf("Error updating state: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Drop complete.")
	},
}

func init() {
	rootCmd.AddCommand(dropCmd)

	// Note: We use the same flag name "target-file" but locally bound to this command
	// We do NOT bind it to viper global config for this command to avoid confusion.
	dropCmd.Flags().String("target-file", "", "Specific target to drop")
}
