package cmd

import (
	"context"
	"fmt"
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
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error")
		return err
	}
	defer file.Close()

	ext := filepath.Ext(path)
	if ext == ".yaml" {
		fmt.Printf("Imported yaml-file: %v\n", file.Name())
		return nil
	} else if ext == ".json" {
		fmt.Printf("Imported json-file: %v\n", file.Name())
		return nil
	}
	return fmt.Errorf("provided file is not a yaml or json file")
}

func init() {
	rootCmd.AddCommand(importCmd)
}
