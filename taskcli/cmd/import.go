package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/MrShanks/Taska/utils"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	apiClient := NewApiClient()

	token := utils.ReadToken()

	result := importTasks(apiClient, ctx, "/import", args[0], token)
	cmd.Printf("%s\n", result)
}

func importTasks(taskcli *Taskcli, ctx context.Context, endpoint, path, token string) string {
	taskcli.ServerURL.Path = endpoint

	file, err := os.Open(path)
	if err != nil {
		return fmt.Sprintf("error opening file %s: %v", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Sprintf("failed to read file %s: %v", path, err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", taskcli.ServerURL.String(), bytes.NewReader(data))
	if err != nil {
		return fmt.Sprintf("failed to create HTTP request: %v", err)
	}

	request.Header.Set("token", token)

	switch ext := filepath.Ext(path); ext {
	case ".yaml":
		request.Header.Set("Content-Type", "application/x-yaml")
	case ".json":
		request.Header.Set("Content-Type", "application/json")
	default:
		return fmt.Sprintln("unsupported file format: only .yaml and .json are allowed")
	}

	response, err := taskcli.HttpClient.Do(request)
	if err != nil {
		return fmt.Sprintf("couldn't get a response from the server: %v", err)
	}
	defer response.Body.Close()

	_, err = CheckStatus(response.StatusCode)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	fmt.Printf("Imported file: %v\n", file.Name())
	return ""
}

func init() {
	rootCmd.AddCommand(importCmd)
}
