package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/utils"
	"github.com/spf13/cobra"
)

var getOneCmd = &cobra.Command{
	Use:   "getone [uuid]",
	Short: "Get one task",
	Long:  "Get one active task from the server",
	Args:  cobra.ExactArgs(1),

	ValidArgsFunction: getCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		apiClient := NewApiClient()

		token := utils.ReadToken()

		output := FetchOne(apiClient, ctx, fmt.Sprintf("/task/%s", args[0]), token)
		outJson, err := json.Marshal(output)
		cobra.CheckErr(err)

		cmd.Printf("%v\n", string(outJson))
	},
}

func getCompletion(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	apiClient := NewApiClient()

	token := utils.ReadToken()

	titles := GetTaskUUIDs(apiClient, ctx, "/tasks", token)

	return titles, cobra.ShellCompDirectiveNoFileComp
}

func GetTaskUUIDs(taskcli *Taskcli, ctx context.Context, endpoint string, token string) []string {
	taskcli.ServerURL.Path = endpoint

	request, err := http.NewRequestWithContext(ctx, "GET", taskcli.ServerURL.String(), nil)
	if err != nil {
		fmt.Printf("Couldn't create request: %v\n", err)
	}

	request.Header.Set("token", token)

	response, err := taskcli.HttpClient.Do(request)
	if err != nil {
		fmt.Printf("Couldn't get a response from the server: %v\n", err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Couldn't read response body: %v\n", err)
	}
	defer response.Body.Close()

	var tasks []*task.Task
	err = json.Unmarshal(bodyBytes, &tasks)
	if err != nil {
		log.Printf("Couldn't decode JSON: %v\n", err)
	}

	uuids := make([]string, len(tasks))
	for _, t := range tasks {
		uuids = append(uuids, t.ID.String())
	}

	return uuids
}

func FetchOne(taskcli *Taskcli, ctx context.Context, endpoint string, token string) *task.Task {
	var t task.Task

	err := fetch(taskcli, ctx, endpoint, &t, token)
	if err != nil {
		log.Printf("Error fetching task: %v\n", err)
		return nil
	}

	return &t
}

func init() {
	rootCmd.AddCommand(getOneCmd)
}
