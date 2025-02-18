package cmd

import (
	"context"
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
		ctx := context.Background()
		httpClient := NewApiClient()

		cmd.Printf("%s", fetchTasks(httpClient, ctx, "/tasks"))
	},
}

func fetchTasks(taskcli *Tasckli, ctx context.Context, endpoint string) string {
	taskcli.ServerURL.Path = endpoint

	request, err := http.NewRequestWithContext(ctx, "GET", taskcli.ServerURL.String(), nil)
	if err != nil {
		log.Printf("Couldn't create request: %v", err)
	}

	response, err := taskcli.HttpClient.Do(request)
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
