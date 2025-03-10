package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var modCmd = &cobra.Command{
	Use:   "mod",
	Short: "Modify a task",
	Long:  "Modify a task by passing its id",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		apiClient := NewApiClient()

		id, err := cmd.Flags().GetString("id")
		if err != nil {
			cmd.Printf("Couldn't retrieve id flag: %v\n", err)
			return
		}

		title, desc, err := getFlags(cmd, 1)
		if err != nil {
			cmd.Printf("Error: %v\n", err)
			return
		}

		body := &bytes.Buffer{}
		body.Write([]byte(fmt.Sprintf(`{"title":"%s","desc":"%s"}`, title, desc)))
		cmd.Printf("%s", modTask(apiClient, ctx, fmt.Sprintf("/update/%s", id), body))
	},
}

func modTask(taskcli *Taskcli, ctx context.Context, endpoint string, body io.Reader) string {
	taskcli.ServerURL.Path = endpoint

	request, err := http.NewRequestWithContext(ctx, "PUT", taskcli.ServerURL.String(), body)
	if err != nil {
		return fmt.Sprintf("Couldn't create request: %v\n", err)
	}

	response, err := taskcli.HttpClient.Do(request)
	if err != nil {
		return fmt.Sprintf("Couldn't get a response from the server: %v\n", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return "Task not found\n"
	}

	return "Task Successfully updated\n"
}

func init() {
	modCmd.Flags().StringP("id", "i", "xxxxx000-0x00-0xx0-x0x0-000x00xx00x0", "Id of the task you want to update")
	modCmd.Flags().StringP("title", "t", "new title", "New title of the task")
	modCmd.Flags().StringP("desc", "d", "new description", "New description of the task")

	err := modCmd.MarkFlagRequired("id")
	cobra.CheckErr(err)

	rootCmd.AddCommand(modCmd)
}
