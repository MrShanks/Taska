package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/MrShanks/Taska/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get all active tasks",
	Long:  "Get or dump on a file all active tasks stored on the server",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		apiClient := NewApiClient()

		format, err := cmd.Flags().GetString("export")
		cobra.CheckErr(err)

		interactive, err := cmd.Flags().GetBool("interactive")
		cobra.CheckErr(err)

		token := utils.ReadToken()

		data := FetchTasks(apiClient, ctx, "/tasks", token)

		if data == nil {
			return
		}

		if interactive {
			items := []list.Item{}

			for _, task := range data {
				items = append(items, item{title: task.Title, desc: task.Desc})
			}

			m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}

			m.list.Styles.Title = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00A5FF")).
				Bold(true).
				Border(lipgloss.RoundedBorder(), true).
				PaddingLeft(3).
				PaddingRight(3)

			m.list.Title = "My Task List"

			p := tea.NewProgram(m, tea.WithAltScreen())

			if _, err := p.Run(); err != nil {
				cmd.Printf("Error running program: %v", err)
				os.Exit(1)
			}
			return
		}

		if isFlagSet(format) {
			var bytes []byte
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

func init() {
	getCmd.Flags().StringP("export", "e", "", "[e]xport tasks in either json|yaml")
	getCmd.Flags().BoolP("interactive", "i", false, "enter [i]nteractive mode")
	getCmd.MarkFlagsMutuallyExclusive("export", "interactive")

	rootCmd.AddCommand(getCmd)
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
