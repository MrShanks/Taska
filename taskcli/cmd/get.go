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
		cmd.Printf("%s", fetchTasks(httpClient, "/tasks"))
	},
}

func fetchTasks(client *http.Client, endpoint string) string {

	serverURL.Path = endpoint

	response, err := client.Get(serverURL.String())
	if err != nil {
		log.Printf("Couldn't get a response from the server: %v", err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Couldn't read response body: %v", err)
	}
	defer response.Body.Close()

	return fmt.Sprintf("%v\n", string(bodyBytes))
}

func init() {
	rootCmd.AddCommand(getCmd)
}
