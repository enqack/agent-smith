package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionDefaults to "dev" if not set via -ldflags
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Agent Smith",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
