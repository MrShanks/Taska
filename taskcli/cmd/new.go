package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/MrShanks/Taska/common/task"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new task",
	Long:  "Create a fancy new task. A task name is required, and a description is optional but recommended",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		apiClient := NewApiClient()

		var title, desc string
		var err error

		if cmd.Flags().NFlag() == 0 {
			title, desc, err = openEditor()
			if err != nil {
				cmd.Printf("couldn't open editor: %v", err)
			}
			if title == "" {
				cmd.Printf("A title must be provided to create a new task")
				return
			}
		} else {
			title, err = cmd.Flags().GetString("title")
			if err != nil {
				log.Printf("Couldn't get the passed value: title")
			}

			desc, err = cmd.Flags().GetString("desc")
			if err != nil {
				log.Printf("Couldn't get the passed value: desc")
			}
		}

		cmd.Printf("%s\n", newTask(apiClient, ctx, "/new", title, desc))
	},
}

// openEditor opens vim editor and returns the edited content.
// First line is interpreted as title and from the second on will be description
func openEditor() (string, string, error) {
	editor := "vim"
	if editor == "" {
		return "", "", fmt.Errorf("no editor found")
	}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "task_edit_*.txt")
	if err != nil {
		return "", "", err
	}
	defer os.Remove(tmpFile.Name()) // Clean up after editing

	// Open the editor
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Wait for the user to edit and save
	if err := cmd.Run(); err != nil {
		return "", "", err
	}

	// Read the edited file content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", "", err
	}

	lines := strings.SplitN(string(content), "\n", 2)
	title := strings.TrimSpace(lines[0]) // First line as title
	description := ""
	if len(lines) > 1 {
		description = strings.TrimSpace(lines[1]) // Rest as description
	}

	return title, description, nil
}

func newTask(taskcli *Tasckli, ctx context.Context, endpoint, title, desc string) string {
	taskcli.ServerURL.Path = endpoint

	jsonTask, err := json.Marshal(task.New(title, desc))
	if err != nil {
		log.Printf("Couldn't marshal task, error: %v", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", taskcli.ServerURL.String(), bytes.NewReader(jsonTask))
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

	// err := newCmd.MarkPersistentFlagRequired("title")
	// if err != nil {
	// 	log.Printf("Error marking persisten flag required: %v", err)
	// }

	rootCmd.AddCommand(newCmd)
}
