package cmd

import (
	"github.com/spf13/cobra"
)

var version = "null"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print taskcli version",
	Long:  `This command reads the version from a config.yaml file`,

	Run: func(cmd *cobra.Command, args []string) {
		if version == "null" {
			version = config.Version
		}
		cmd.Printf("version: %s", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
