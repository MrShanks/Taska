package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// openEditor opens vim editor and returns the edited content.
// First line is interpreted as title and from the second on will be description
func openEditor() (string, string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "task_edit_*.txt")
	if err != nil {
		return "", "", err
	}
	defer os.Remove(tmpFile.Name()) // Clean up after editing

	// Open the editor
	cmd := exec.Command("vim", tmpFile.Name()) //nolint:gosec // input is passed directly
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

func getFlags(cmd *cobra.Command, required int) (string, string, error) {
	title, err := cmd.Flags().GetString("title")
	if err != nil {
		return "", "", fmt.Errorf("couldn't get the passed value: title: %v", err)
	}

	desc, err := cmd.Flags().GetString("desc")
	if err != nil {
		return title, "", fmt.Errorf("couldn't get the passed value: desc: %v", err)
	}

	if cmd.Flags().NFlag() == required {
		// No title or desc provided, open the editor
		title, desc, err = openEditor()
		if err != nil {
			return "", "", fmt.Errorf("couldn't open editor: %v", err)
		}
	}

	if title == "" && desc == "" {
		return "", "", fmt.Errorf("you must provide at least a title or a desc flag to modify a task")

	}

	return title, desc, nil
}

func fetch(taskcli *Taskcli, ctx context.Context, endpoint string, result any) error {
	taskcli.ServerURL.Path = endpoint

	request, err := http.NewRequestWithContext(ctx, "GET", taskcli.ServerURL.String(), nil)
	if err != nil {
		return fmt.Errorf("Couldn't create request: %v", err)
	}

	response, err := taskcli.HttpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Couldn't get a response from the server: %v", err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Couldn't read response body: %v", err)
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		return fmt.Errorf("Couldn't unmarshal: %v", err)
	}

	return nil
}
