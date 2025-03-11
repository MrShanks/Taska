package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/MrShanks/Taska/common/task"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var getOneCmd = &cobra.Command{
	Use:   "getone [uuid]",
	Short: "Get one task",
	Long:  "Get one active task from the server",
	Args:  cobra.ExactArgs(1),

	ValidArgsFunction: getCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		apiClient := NewApiClient()
		ctx := context.Background()

		output := FetchOne(apiClient, ctx, fmt.Sprintf("/task/%s", args[0]))
		outJson, err := json.Marshal(output)
		cobra.CheckErr(err)

		cmd.Printf("%v\n", string(outJson))
	},
}

func getCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	apiClient := NewApiClient()
	ctx := context.Background()
	titles := GetTaskUUIDs(apiClient, ctx, "/tasks")

	return titles, cobra.ShellCompDirectiveNoFileComp
}

func GetTaskUUIDs(taskcli *Taskcli, ctx context.Context, endpoint string) []string {
	taskcli.ServerURL.Path = endpoint

	request, err := http.NewRequestWithContext(ctx, "GET", taskcli.ServerURL.String(), nil)
	if err != nil {
		fmt.Printf("Couldn't create request: %v", err)
	}

	response, err := taskcli.HttpClient.Do(request)
	if err != nil {
		fmt.Printf("Couldn't get a response from the server: %v", err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Couldn't read response body: %v", err)
	}
	defer response.Body.Close()

	var tasks map[uuid.UUID]*task.Task
	err = json.Unmarshal(bodyBytes, &tasks)
	if err != nil {
		log.Printf("Couldn't decode JSON: %v", err)
	}

	uuids := make([]string, len(tasks))
	for _, t := range tasks {
		uuids = append(uuids, t.ID.String())
	}

	return uuids
}

func FetchOne(taskcli *Taskcli, ctx context.Context, endpoint string) *task.Task {
	var t task.Task

	err := fetch(taskcli, ctx, endpoint, &t)
	if err != nil {
		log.Printf("Error fetching task: %v", err)
		return nil
	}

	return &t
}

func init() {
	rootCmd.AddCommand(getOneCmd)
}
