package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"

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
			log.Printf("Couldn't get a response from the server: %v", err)
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Couldn't read response body: %v", err)
		}

		fmt.Printf("%v\n", string(bodyBytes))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
