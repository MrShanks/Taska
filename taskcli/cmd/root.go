package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "taskcli",
	Short: "taskcli is your CLI best friend",
	Long: `taskcli is your CLI best friend. 
	
handling your tasks has never been easier,
you can create, modify, delete and read your tasks

Example usage:
  taskcli get 
  taskcli new -t <task_title> -d <task_desc>
  taskcli mod -i <task_id> -t <new_title> -d <new_desc>
  taskcli del -i <task_id>`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// add flags here
}
