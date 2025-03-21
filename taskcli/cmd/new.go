package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/utils"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new task",
	Long:  "Create a fancy new task. A task name is required, and a description is optional but recommended",

	Run: runNewCmd,
}

func runNewCmd(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	apiClient := NewApiClient()

	title, desc, err := getFlags(cmd, 0)
	cobra.CheckErr(err)

	if title == "" {
		cmd.Printf("A title must be provided to create a new task\n")
		return
	}

	token := utils.ReadToken()

	result := newTask(apiClient, ctx, "/new", title, desc, token)
	cmd.Printf("%s\n", result)
}

func newTask(taskcli *Taskcli, ctx context.Context, endpoint, title, desc, token string) string {
	taskcli.ServerURL.Path = endpoint

	jsonTask, err := json.Marshal(task.New(title, desc))
	if err != nil {
		return fmt.Sprintf("Couldn't marshal task, error: %v", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", taskcli.ServerURL.String(), bytes.NewReader(jsonTask))
	if err != nil {
		return fmt.Sprintf("Couldn't create request: %v", err)
	}

	request.Header.Set("token", token)

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
	newCmd.Flags().StringP("title", "t", "Untitled task", "Title of the new fancy task")
	newCmd.Flags().StringP("desc", "d", "Default description", "Description of the new task")

	rootCmd.AddCommand(newCmd)
}
