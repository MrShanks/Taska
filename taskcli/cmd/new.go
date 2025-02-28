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

	Run: runNewCmd,
}

func runNewCmd(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	apiClient := NewApiClient()

	title, desc, err := getFlags(cmd, 0)
	if err != nil {
		cmd.Printf("Error: %v\n", err)
		return
	}

	if title == "" {
		cmd.Printf("A title must be provided to create a new task\n")
		return
	}

	result := newTask(apiClient, ctx, "/new", title, desc)
	cmd.Printf("%s\n", result)
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

	rootCmd.AddCommand(newCmd)
}
