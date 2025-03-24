package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/MrShanks/Taska/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search a task by keyword",
	Long: `Search a task by keyword, all matching tasks that have the keyword
either in the title or in the description will be printed on the stdout`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		apiClient := NewApiClient()
		token := utils.ReadToken()
		keyword := args[0]
		endpoint := fmt.Sprintf("/search/%s", keyword)

		output := search(apiClient, ctx, endpoint, token)

		cmd.Printf("Found: %s\n", color.YellowString(output))
	},
}

func search(taskcli *Taskcli, ctx context.Context, endpoint, token string) string {
	return endpoint
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
