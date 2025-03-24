package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/MrShanks/Taska/utils"
	"github.com/spf13/cobra"
)

var (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
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

		cmd.Printf("Found: "+Yellow+"%s\n"+Reset, output)
	},
}

func search(taskcli *Taskcli, ctx context.Context, endpoint, token string) string {
	return endpoint
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
