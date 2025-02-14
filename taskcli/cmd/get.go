package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MrShanks/Taska/common/logger"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get all active tasks",
	Long:  "get all active tasks store on the servere",

	Run: func(cmd *cobra.Command, args []string) {
		var apiClient http.Client

		response, err := apiClient.Get("http://localhost:8080/tasks")
		if err != nil {
			logger.ErrorLogger.Printf("Couldn't get a response from the server: %v", err)
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			logger.ErrorLogger.Printf("Couldn't read response body: %v", err)
		}

		fmt.Printf("%v\n", string(bodyBytes))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
