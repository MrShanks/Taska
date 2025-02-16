package cmd

import (
	"net/http"
	"net/url"
	"os"

	"github.com/MrShanks/Taska/utils"
	"github.com/spf13/cobra"
)

var (
	httpClient *http.Client
	config     *utils.Config
	serverURL  url.URL
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "taskcli",
	Short: "taskcli is your CLI best friend",
	Long: `taskcli is your CLI best friend. 
	
handling your tasks has never been easier,
you can create, modify, delete and read your tasks

Example usage:
  taskcli new <task_title> <task_desc>
  taskcli mod <task_id> -t <new_title> -d <new_desc>
  taskcli get <task_id>
  taskcli del <task_id>`,
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
