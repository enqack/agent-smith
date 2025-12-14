package cli

import (
	"agent-smith/internal/config"
	"agent-smith/internal/ops"
	"agent-smith/internal/state"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current persona status",
	Long:  `Show which persona is currently active by checking the AGENTS.md symlink.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load state first
		st, err := state.LoadState()
		if err != nil {
			// Proceed with empty state if loading fails (first run?)
			st = &state.StatusState{}
		}

		// Determine Canonical Target
		canonical := st.CanonicalTarget
		if canonical == "" {
			canonical = viper.GetString("target_file")
		}

		// Infer active persona
		activePersona := inferPersona(canonical)

		fmt.Printf("Status Check:\n\n")

		if len(st.AgentFiles) == 0 {
			// Check if we have active persona even without state (legacy/fresh)
			if activePersona != "" {
				fmt.Printf("Active Persona: %s (Not tracked in state)\n", activePersona)
				fmt.Println(" Targets (from config):")
				for _, t := range Cfg.Targets {
					printTargetStatus(state.TargetState{Path: t.Path, Mode: t.Mode}, activePersona)
				}
			} else {
				fmt.Println("No active persona and no state found.")
			}
			return
		}

		// Iterate over all known personas in state
		for _, af := range st.AgentFiles {
			isActive := (af.Name == activePersona)
			statusStr := ""
			if isActive {
				statusStr = " [ACTIVE]"
			}

			fmt.Printf("Persona: %s%s\n", af.Name, statusStr)
			if af.Path != "" {
				fmt.Printf("  File: %s\n", af.Path)
			}
			fmt.Println("  Targets:")

			for _, t := range af.Targets {
				printTargetStatus(t, af.Name)
			}
			fmt.Println()
		}

		// If there is an active persona NOT in state?
		foundActiveInState := false
		for _, af := range st.AgentFiles {
			if af.Name == activePersona {
				foundActiveInState = true
				break
			}
		}

		if activePersona != "" && !foundActiveInState {
			fmt.Printf("Persona: %s [ACTIVE] (Config only)\n", activePersona)
			fmt.Println("  Targets:")
			for _, t := range Cfg.Targets {
				printTargetStatus(state.TargetState{Path: t.Path, Mode: t.Mode}, activePersona)
			}
		}
	},
}

func inferPersona(path string) string {
	path = ops.ExpandPath(path)

	info, err := os.Lstat(path)
	if err != nil {
		return ""
	}

	if info.Mode()&os.ModeSymlink != 0 {
		dest, err := os.Readlink(path)
		if err != nil {
			return ""
		}
		// Expected: .../AGENTS.<persona>.md
		base := filepath.Base(dest)
		if len(base) > 10 && base[:7] == "AGENTS." && base[len(base)-3:] == ".md" {
			return base[7 : len(base)-3]
		}
	}
	return ""
}

func printTargetStatus(target state.TargetState, personaName string) {
	targetPath := ops.ExpandPath(target.Path)
	// We need the source path to verify
	DisplayPath := targetPath
	if len(targetPath) > 40 {
		DisplayPath = "..." + targetPath[len(targetPath)-37:]
	}

	info, err := os.Lstat(targetPath)

	status := "OK"
	details := ""

	if err != nil {
		if os.IsNotExist(err) {
			status = "MISSING"
		} else {
			status = "ERROR"
			details = fmt.Sprintf("(%v)", err)
		}
	} else {
		if target.Mode == config.TargetModeLink {
			if info.Mode()&os.ModeSymlink == 0 {
				status = "DRIFT"
				details = "(Not a symlink)"
			} else {
				linkDest, err := os.Readlink(targetPath)
				if err != nil {
					status = "ERROR"
					details = fmt.Sprintf("(Readlink failed: %v)", err)
				} else {
					expectedSuffix := fmt.Sprintf("AGENTS.%s.md", personaName)
					if filepath.Base(linkDest) != expectedSuffix {
						status = "DRIFT"
						details = fmt.Sprintf("(Points to %s)", filepath.Base(linkDest))
					}
				}
			}
		} else if target.Mode == config.TargetModeCopy {
			if info.Mode()&os.ModeSymlink != 0 {
				status = "DRIFT"
				details = "(Is a symlink, expected copy)"
			} else {
				// Copy check - maybe check content hash? Too expensive?
				// Just checking existence for now is fine for "OK" vs "MISSING".
				// "Drift" for copy implies content mismatch, which we can't easily check without source.
			}
		}
	}

	fmt.Printf("    [%s] %s %s\n", status, DisplayPath, details)
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
