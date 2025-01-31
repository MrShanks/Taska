package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "undefined"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print taskcli version",
	Long:  `This command reads the version from ~/.Taska.yaml config file`,

	Run: func(cmd *cobra.Command, args []string) {
		if version == "undefined" {
			version = readVersionFromConfig()
		}
		fmt.Printf("version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
