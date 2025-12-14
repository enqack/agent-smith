package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/agent-smith/config.yaml)")
	rootCmd.PersistentFlags().StringSlice("agents-dir", []string{}, "directory containing agent personas (can be specified multiple times)")
	rootCmd.PersistentFlags().String("target-file", "", "path to the AGENTS.md symlink")

	// Bind flags to viper
	viper.BindPFlag("agents_dir", rootCmd.PersistentFlags().Lookup("agents-dir"))
	viper.BindPFlag("target_file", rootCmd.PersistentFlags().Lookup("target-file"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in XDG_CONFIG_HOME/agent-smith (e.g. ~/.config/agent-smith)
		// and /etc/agent-smith
		configHome, err := getConfigHome()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(filepath.Join(configHome, "agent-smith"))
		viper.AddConfigPath("/etc/agent-smith")
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
	dataHome, err := getDataHome()
	if err == nil {
		defaultAgentsDirs = append(defaultAgentsDirs, filepath.Join(dataHome, "agent-smith", "personas"))
	}
	defaultAgentsDirs = append(defaultAgentsDirs, "/usr/share/agent-smith/personas")

	viper.SetDefault("agents_dir", defaultAgentsDirs)

	configHome, err := getConfigHome()
	if err == nil {
		viper.SetDefault("target_file", filepath.Join(configHome, "agents", "AGENTS.md"))
	} else {
		// Fallback
		viper.SetDefault("target_file", "AGENTS.md")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
