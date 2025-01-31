package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print taskcli version",
	Long:  `This command reads the version from ~/.Taska.yaml config file`,

	Run: func(cmd *cobra.Command, args []string) {
		version := viper.GetString("version")

		if version == "" {
			fmt.Fprintln(os.Stderr, "version not set in ~/.Taska.yaml")
			os.Exit(1)
		}

		fmt.Println("version: ", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
