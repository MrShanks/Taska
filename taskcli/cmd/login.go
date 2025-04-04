package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MrShanks/Taska/utils"
	"github.com/spf13/cobra"
)

var token string

var loginCmd = &cobra.Command{
	Use:   "login -e [email] -p [password]",
	Short: "Login with email and password",
	Long:  "Login with email and password to get a token to auhtenticate future requests",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		apiClient := NewApiClient()

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

	response, err := makeRequest(taskcli, ctx, data)
	if err != nil {
		return fmt.Errorf("error while making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%v", response.StatusCode)
	}

	token = response.Header.Get("Token")

	utils.StoreToken(token)

	srvMsg, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("couldn't read response body: %v", err)
	}

	fmt.Printf("%v\n", string(srvMsg))

	return nil
}

func init() {
	loginCmd.Flags().StringP("email", "e", "marco@rossi.com", "email to authenticate author")
	loginCmd.Flags().StringP("password", "p", "password", "password to authenticate author")

	err := loginCmd.MarkFlagRequired("email")
	cobra.CheckErr(err)

	err = loginCmd.MarkFlagRequired("password")
	cobra.CheckErr(err)

	rootCmd.AddCommand(loginCmd)
}
