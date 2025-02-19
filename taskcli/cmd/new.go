package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/MrShanks/Taska/common/task"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create a new task",
	Long:  "create a fancy new task. A task name is required, and a description is optional but recommended",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		httpClient := NewApiClient()
		title, _ := cmd.Flags().GetString("title")
		desc, _ := cmd.Flags().GetString("desc")

		cmd.Printf("%s", newTask(httpClient, ctx, "/new", title, desc))
	},
}

func newTask(taskcli *Tasckli, ctx context.Context, endpoint, title, desc string) string {
	taskcli.ServerURL.Path = endpoint

	bytesNewTask, err := json.Marshal(task.New(title, desc))
	if err != nil {
		log.Printf("Couldn't marshal received task: %v", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", taskcli.ServerURL.String(), bytes.NewReader(bytesNewTask))
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
	newCmd.PersistentFlags().StringP("title", "t", "Untitled task", "Title of the new fancy task")
	newCmd.PersistentFlags().StringP("desc", "d", "Default description", "Description of the new task")
	err := newCmd.MarkPersistentFlagRequired("title")
	if err != nil {
		log.Printf("Error marking persisten flag required: %v", err)
	}

	rootCmd.AddCommand(newCmd)
}
