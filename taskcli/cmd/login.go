package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var token string

var loginCmd = &cobra.Command{
	Use:   "login -e [email] -p [password]",
	Short: "Login with email and password",
	Long:  "Login with email and password to get a token to auhtenticate future requests",

	Run: func(cmd *cobra.Command, args []string) {
		apiClient := NewApiClient()
		ctx := context.Background()

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			cmd.Printf("A email must be provided to login: %v\n", err)
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			cmd.Printf("Empty password is not allowed: %v\n", err)
		}

		if err := Login(apiClient, ctx, "signin", email, password); err != nil {
			cmd.Printf("Login failed: %v\n", err)
		}
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
		return fmt.Errorf("%v", response.StatusCode)
	}

	token = response.Header.Get("Token")

	storeToken(token)

	srvMsg, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Couldn't read response body: %v", err)
	}

	fmt.Printf("%v\n", string(srvMsg))

	return nil
}

func storeToken(token string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("couldn't locate user home")
	}

	file, err := os.Create(filepath.Join(home, ".taskcli"))
	if err != nil {
		fmt.Printf("couldn't store login credentials: %v", err)
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n", token))
	if err != nil {
		fmt.Printf("cannot write to %s, error: %v", file.Name(), err)
	}

}

func init() {
	loginCmd.Flags().StringP("email", "e", "marco@rossi.com", "email to authenticate author")
	loginCmd.Flags().StringP("password", "p", "password", "password to authenticate author")

	rootCmd.AddCommand(loginCmd)
}
