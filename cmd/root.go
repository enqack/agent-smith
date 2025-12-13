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
	Long: `Agents is a CLI tool for managing AGENTS.md symlinks,
allowing users to switch between different agent personas.`,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/agents/config.yaml)")
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
		// Search config in home directory and /etc/agents
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(filepath.Join(home, ".config", "agents"))
		viper.AddConfigPath("/etc/agents")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		// Ensure local config directory exists
		configDir := filepath.Join(home, ".config", "agents")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			// Not fatal, but good to know
			// fmt.Printf("Warning: Could not create config directory: %v\n", err)
		}
	}

	// Set defaults
	var defaultAgentsDirs []string
	home, err := os.UserHomeDir()
	if err == nil {
		defaultAgentsDirs = append(defaultAgentsDirs, filepath.Join(home, ".config", "agents"))
		defaultAgentsDirs = append(defaultAgentsDirs, filepath.Join(home, ".config", "agent-smith", "agents"))
	}
	defaultAgentsDirs = append(defaultAgentsDirs, "/usr/share/agent-smith/agents")

	viper.SetDefault("agents_dir", defaultAgentsDirs)
	if err == nil {
		viper.SetDefault("target_file", filepath.Join(home, ".config", "agents", "AGENTS.md"))
	} else {
		// Fallback if home directory not found (unlikely for user)
		viper.SetDefault("target_file", "AGENTS.md")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
