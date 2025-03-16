package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login -e [email] -p [password]",
	Short: "Login with email and password",
	Long:  "Login with email and password to get a token to auhtenticate future requests",

	Run: func(cmd *cobra.Command, args []string) {
		apiClient := NewApiClient()
		ctx := context.Background()

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			cmd.Printf("something went wrong when fetching the email: %v", err)
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			cmd.Printf("something went wrong when fetching the email: %v", err)
		}

		Login(apiClient, ctx, "login", email, password)
	},
}

func Login(taskcli *Taskcli, ctx context.Context, endpoint, email, password string) error {
	taskcli.ServerURL.Path = endpoint

	data := map[string]string{"email": email, "password": password}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Couldn't marshal data into the request body: %v", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", taskcli.ServerURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Couldn't create request: %v", err)
	}

	response, err := taskcli.HttpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Couldn't get a response from the server: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("task not found, status code: %v", response.StatusCode)
	}
	// TODO: FInish this one up

	return nil

}

func init() {
	newCmd.Flags().StringP("email", "e", "marco@rossi.com", "email to authenticate author")
	newCmd.Flags().StringP("password", "p", "password", "password to authenticate author")

	rootCmd.AddCommand(loginCmd)
}
