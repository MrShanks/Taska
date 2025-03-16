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
		cmd.Printf("no valid path, error: %v", err)
		return
	}
}

func importTasks(taskcli *Taskcli, ctx context.Context, endpoint, path string) error {
	taskcli.ServerURL.Path = endpoint

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error")
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
	}

	ext := filepath.Ext(path)
	if ext == ".yaml" || ext == ".yml" {
		request, err := http.NewRequestWithContext(ctx, "POST", taskcli.ServerURL.String(), bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("couldn't create request: %v", err)
		}

		request.Header.Set("Content-Type", "application/x-yaml")

		_, err = taskcli.HttpClient.Do(request)
		if err != nil {
			return fmt.Errorf("couldn't get a response from the server: %v", err)
		}

		fmt.Printf("Imported yaml-file: %v\n", file.Name())
		return nil
	} else if ext == ".json" {
		request, err := http.NewRequestWithContext(ctx, "POST", taskcli.ServerURL.String(), bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("couldn't create request: %v", err)
		}

		request.Header.Set("Content-Type", "application/json")

		_, err = taskcli.HttpClient.Do(request)
		if err != nil {
			return fmt.Errorf("couldn't get a response from the server: %v", err)
		}

		fmt.Printf("Imported json-file: %v\n", file.Name())
		return nil
	}
	return fmt.Errorf("provided file is not a yaml or json file")
}

func init() {
	rootCmd.AddCommand(importCmd)
}
