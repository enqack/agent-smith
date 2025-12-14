package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"agent-smith/internal/config"
	"agent-smith/internal/ops"
)

var (
	// Used for flags.
	cfgFile string

	// Cfg stores the global configuration
	Cfg config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "agents",
	Short: "An agent version control program",
	Long: `agents - Agent Smith persona manager

Manage AGENTS.md symlinks to switch between different agent personas.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/agent-smith/config.yaml)")

	// Bind flags to viper (legacy flags kept if needed, but per previous edits we want standard XDG)
	// We previously had agents-dir etc. Let's keep them reachable via flags if user wants overrides.
	rootCmd.PersistentFlags().StringSlice("agents-dir", []string{}, "directory containing agent personas (can be specified multiple times)")
	rootCmd.PersistentFlags().String("target-file", "", "path to the AGENTS.md symlink")

	// Bind agents_dir to viper
	viper.BindPFlag("agents_dir", rootCmd.PersistentFlags().Lookup("agents-dir"))

	// Note: We DO NOT bind "target-file" flag to "target_file" config.
	// The flag is ephemeral (where to write NOW), the config is persistent (System Canonical Path).
	// viper.BindPFlag("target_file", rootCmd.PersistentFlags().Lookup("target-file"))

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find config directory
		configHome, err := config.GetConfigHome()
		if err != nil {
			fmt.Println("Error getting config home:", err)
			os.Exit(1)
		}

		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(filepath.Join(configHome, "agent-smith"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		// Ensure local config directory exists
		configDir := filepath.Join(configHome, "agent-smith")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			// Not fatal, but good to know
			// fmt.Printf("Warning: Could not create config directory: %v\n", err)
		}
	}

	// Set defaults
	var defaultAgentsDirs []string
	dataHome, err := config.GetDataHome()
	if err == nil {
		defaultAgentsDirs = append(defaultAgentsDirs, filepath.Join(dataHome, "agent-smith", "personas"))
	}
	defaultAgentsDirs = append(defaultAgentsDirs, "/usr/share/agent-smith/personas")

	viper.SetDefault("agents_dir", defaultAgentsDirs)

	cHome, err := config.GetConfigHome()
	if err == nil {
		// Canonical Target: $XDG_CONFIG_HOME/agents/AGENTS.md
		viper.SetDefault("target_file", filepath.Join(cHome, "agents", "AGENTS.md"))
	} else {
		// Fallback
		viper.SetDefault("target_file", "AGENTS.md")
	}

	viper.SetEnvPrefix("AGENTS")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Unmarshal config
	if err := viper.Unmarshal(&Cfg); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	// Backward Compatibility:
	// If 'target_file' is set but not in 'targets', add it as a managed target (LINK mode).
	// This ensures legacy users still get their main symlink managed/monitored.
	if Cfg.TargetFile != "" {
		found := false
		absTarget, _ := filepath.Abs(ops.ExpandPath(Cfg.TargetFile))

		for _, t := range Cfg.Targets {
			tAbs, _ := filepath.Abs(ops.ExpandPath(t.Path))
			if tAbs == absTarget {
				found = true
				break
			}
		}
		if !found {
			Cfg.Targets = append(Cfg.Targets, config.TargetConfig{
				Path: Cfg.TargetFile,
				Mode: config.TargetModeLink,
			})
		}
	}
}
