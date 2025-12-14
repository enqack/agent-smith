package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"agent-smith/internal/config"
	"agent-smith/internal/ops"
	"agent-smith/internal/state"
)

// reconcileCmd represents the reconcile command
var reconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reapply the active persona to all targets",
	Long:  `Reapply the currently active persona to all configured targets, fixing any drift or missing files.`,
	Run: func(cmd *cobra.Command, args []string) {
		st, err := state.LoadState()
		if err != nil || st == nil || len(st.AgentFiles) == 0 {
			fmt.Println("No active personas (agent files) found. Cannot reconcile.")
			return
		}

		agentsDirs := viper.GetStringSlice("agents_dir")
		if len(agentsDirs) == 0 {
			if s := viper.GetString("agents_dir"); s != "" {
				agentsDirs = []string{s}
			}
		}

		// Determine Canonical Target and Active Persona
		canonical := st.CanonicalTarget
		if canonical == "" {
			canonical = viper.GetString("target_file")
		}

		activePersona := inferPersona(canonical)
		if activePersona == "" {
			fmt.Println("No active persona found to reconcile.")
			return
		}

		fmt.Printf("Reconciling active persona: %s\n", activePersona)

		// Find the active persona in state to get its tracked targets
		// This ensures we respect dynamic targets (CLI flags) that were saved.
		var targetsToApply []config.TargetConfig
		foundInState := false

		for _, af := range st.AgentFiles {
			if af.Name == activePersona {
				for _, t := range af.Targets {
					targetsToApply = append(targetsToApply, config.TargetConfig{
						Path: t.Path,
						Mode: t.Mode,
					})
				}
				foundInState = true
				break
			}
		}

		if !foundInState {
			// Fallback to Config targets if not tracked in state yet (legacy or manual switch)
			// But note: if it WAS in state but had 0 targets, targetsToApply is empty, which is correct.
			// But we only set foundInState if we found the entry.
			// If not found, use config.
			fmt.Println("Active persona not found in state, using defaults.")
			targetsToApply = Cfg.Targets
		}

		// Reapply active persona to targets
		_, err = ops.ApplyPersona(activePersona, agentsDirs, targetsToApply)
		if err != nil {
			fmt.Printf("Failed to reconcile: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Reconciliation complete.")
	},
}

func init() {
	rootCmd.AddCommand(reconcileCmd)
}
