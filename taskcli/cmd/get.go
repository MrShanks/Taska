package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MrShanks/Taska/common/task"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get all active tasks",
	Long:  "Get or dump on a file all active tasks stored on the server",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		apiClient := NewApiClient()

		format, err := cmd.Flags().GetString("export")
		cobra.CheckErr(err)

		token := readToken()

		data := FetchTasks(apiClient, ctx, "/tasks", token)

		if data == nil {
			cmd.Printf("Server response was empty\n")
			return
		}

		var bytes []byte

		if isFlagSet(format) {
			switch format {
			case "yaml":
				bytes, err = yaml.Marshal(data)
				if err != nil {
					cmd.Printf("Error during yaml marshalling: %v\n", err)
				}
				dumpOnFile("export", format, bytes)
			case "json":
				bytes, err = json.MarshalIndent(data, "", "	")
				if err != nil {
					cmd.Printf("Couldn't marshal indent: %v\n", err)
				}
				dumpOnFile("export", format, bytes)
			default:
				cmd.Printf("Unsupported file format: %s choose between json|yaml\n", format)
			}
			return
		}

		output, err := json.Marshal(data)
		cobra.CheckErr(err)

		cmd.Printf("%s\n", string(output))
	},
}

func isFlagSet(flagValue string) bool {
	return flagValue != ""
}

func dumpOnFile(filepath, format string, data []byte) {
	file, err := os.Create(fmt.Sprintf("%s.%s", filepath, format))
	if err != nil {
		fmt.Printf("Couldn't create an export file: %v", err)
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n", data))
	if err != nil {
		fmt.Printf("cannot write to %s, error: %v\n", file.Name(), err)
	}

	fmt.Printf("file created: %s.%s\n", "export", format)
}

func FetchTasks(taskcli *Taskcli, ctx context.Context, endpoint string, token string) []*task.Task {
	var tasks []*task.Task

	err := fetch(taskcli, ctx, endpoint, &tasks, token)
	if err != nil {
		fmt.Printf("Error fetching tasks: %v\n", err)
		return nil
	}

	return tasks
}

func init() {
	getCmd.Flags().StringP("export", "e", "", "export tasks in either json|yaml")

	rootCmd.AddCommand(getCmd)
}
