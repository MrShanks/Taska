package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MrShanks/Taska/utils"
	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a task",
	Long:  "Delete a task by passing its id",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		apiClient := NewApiClient()

		token := utils.ReadToken()

		id, err := cmd.Flags().GetString("id")
		if err != nil {
			cmd.Printf("Couldn't retrieve id flag: %v", err)
		}

		cmd.Printf("%s", delTask(apiClient, ctx, fmt.Sprintf("/delete/%s", id), token))
	},
}

func delTask(taskcli *Taskcli, ctx context.Context, endpoint, token string) string {
	taskcli.ServerURL.Path = endpoint

	request, err := http.NewRequestWithContext(ctx, "DELETE", taskcli.ServerURL.String(), nil)
	if err != nil {
		return fmt.Sprintf("Couldn't create request: %v", err)
	}

	request.Header.Set("token", token)

	response, err := taskcli.HttpClient.Do(request)
	if err != nil {
		return fmt.Sprintf("Couldn't get a response from the server: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return "Task not found"
	}

	if response.StatusCode == http.StatusBadRequest {
		return "Uuid not valid"
	}

	return "Task Successfully deleted"
}

func init() {
	delCmd.PersistentFlags().StringP("id", "i", "xxxxx000-0x00-0xx0-x0x0-000x00xx00x0", "Id of the task you want to delete")

	err := delCmd.MarkPersistentFlagRequired("id")
	if err != nil {
		log.Printf("Error marking persisten flag required: %v", err)
	}

	rootCmd.AddCommand(delCmd)
}
