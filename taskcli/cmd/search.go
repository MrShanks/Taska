package cmd

import (
	"context"
	"fmt"
	"strings"
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

		data := FetchTasks(apiClient, ctx, endpoint, token)

		for _, task := range data {
			fmt.Printf("%-10s %v\n", "ID:", task.ID)
			fmt.Printf("%-10s %s\n", "Title:", HighlightText(task.Title, keyword))
			fmt.Printf("%-10s %s\n", "Desc:", HighlightText(task.Desc, keyword))
			fmt.Printf("%-10s %v\n\n", "Author:", task.AuthorID)
		}
	},
}

func HighlightText(text, keyword string) string {
	if keyword == "" {
		return text
	}

	highlightedKeyword := color.YellowString(keyword)
	return strings.ReplaceAll(text, keyword, highlightedKeyword)
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
