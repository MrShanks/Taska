package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/MrShanks/Taska/common/task"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new task",
	Long:  "Create a fancy new task. A task name is required, and a description is optional but recommended",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		apiClient := NewApiClient()
		title, err := cmd.Flags().GetString("title")
		if err != nil {
			cmd.Printf("Couldn't get the passed value: title")
		}
		desc, err := cmd.Flags().GetString("desc")
		if err != nil {
			cmd.Printf("Couldn't get the passed value: desc")
		}

		cmd.Printf("%s", newTask(apiClient, ctx, "/new", title, desc))
	},
}

func newTask(taskcli *Taskcli, ctx context.Context, endpoint, title, desc string) string {
	taskcli.ServerURL.Path = endpoint

	jsonTask, err := json.Marshal(task.New(title, desc))
	if err != nil {
		return fmt.Sprintf("Couldn't marshal task, error: %v", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", taskcli.ServerURL.String(), bytes.NewReader(jsonTask))
	if err != nil {
		return fmt.Sprintf("Couldn't create request: %v", err)
	}

	response, err := taskcli.HttpClient.Do(request)
	if err != nil {
		return fmt.Sprintf("Couldn't get a response from the server: %v", err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Sprintf("Couldn't read response body: %v", err)
	}
	defer response.Body.Close()

	return fmt.Sprintf("%v\n", string(bodyBytes))
}

func init() {
	newCmd.PersistentFlags().StringP("title", "t", "Untitled task", "Title of the new fancy task")
	newCmd.PersistentFlags().StringP("desc", "d", "Default description", "Description of the new task")

	err := newCmd.MarkPersistentFlagRequired("title")
	if err != nil {
		newCmd.Printf("Error marking persisten flag required: %v", err)
	}

	rootCmd.AddCommand(newCmd)
}
