package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import tasks",
	Long:  "Import tasks from yaml or json file",
	Args:  cobra.ExactArgs(1),

	Run: runImportCmd,
}

func runImportCmd(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	apiClient := NewApiClient()

	err := importTasks(apiClient, ctx, "/import", args[0])
	if err != nil {
		cmd.Printf("error running importTasks: %v", err)
		return
	}
}

func importTasks(taskcli *Taskcli, ctx context.Context, endpoint, path string) error {
	taskcli.ServerURL.Path = endpoint

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", taskcli.ServerURL.String(), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	switch ext := filepath.Ext(path); ext {
	case ".yaml":
		request.Header.Set("Content-Type", "application/x-yaml")
	case ".json":
		request.Header.Set("Content-Type", "application/json")
	default:
		return fmt.Errorf("unsupported file format: only .yaml and .json are allowed")
	}

	response, err := taskcli.HttpClient.Do(request)
	if err != nil {
		return fmt.Errorf("couldn't get a response from the server: %v", err)
	}
	defer response.Body.Close()
	fmt.Printf("Imported file: %v\n", file.Name())
	return nil
}

func init() {
	rootCmd.AddCommand(importCmd)
}
