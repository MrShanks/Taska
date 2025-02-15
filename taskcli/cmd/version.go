package cmd

import (
	"fmt"

	"github.com/MrShanks/Taska/utils"
	"github.com/spf13/cobra"
)

var version = "undefined"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print taskcli version",
	Long:  `This command reads the version from a config.yaml file`,

	Run: func(cmd *cobra.Command, args []string) {
		if version == "undefined" {
			version = utils.LoadConfig().Version
		}
		fmt.Printf("version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
