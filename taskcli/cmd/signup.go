package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MrShanks/Taska/common/author"
	"github.com/spf13/cobra"
)

var signupCmd = &cobra.Command{
	Use:   "signup -f [firstname] -l [lastname] -e [email] -p [password]",
	Short: "Signup a new author",
	Long:  "Signup a new auhtor to be able to login",

	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		apiClient := NewApiClient()

		firstname, err := cmd.Flags().GetString("firstname")
		cobra.CheckErr(err)

		lastname, err := cmd.Flags().GetString("lastname")
		cobra.CheckErr(err)

		email, err := cmd.Flags().GetString("email")
		cobra.CheckErr(err)

		password, err := cmd.Flags().GetString("password")
		cobra.CheckErr(err)

		newAuthor := author.Author{
			Firstname: firstname,
			Lastname:  lastname,
			Email:     email,
			Password:  password,
		}

		if err := Signup(apiClient, ctx, "signup", &newAuthor); err != nil {
			cmd.Printf("Signup failed: %v\n", err)
		}
	},
}

func Signup(taskcli *Taskcli, ctx context.Context, endpoint string, author *author.Author) error {
	taskcli.ServerURL.Path = endpoint

	data := map[string]string{"firstname": author.Firstname, "lastname": author.Lastname, "email": author.Email, "password": author.Password}

	response, err := makeRequest(taskcli, ctx, data)
	if err != nil {
		return fmt.Errorf("error while making request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("%v", response.StatusCode)
	}

	srvMsg, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("couldn't read response body: %v", err)
	}

	fmt.Printf("%v\n", string(srvMsg))

	return nil
}

func init() {
	signupCmd.Flags().StringP("firstname", "f", "marco", "firstname of the author")
	signupCmd.Flags().StringP("lastname", "l", "rossi", "lastname of the author")

	signupCmd.Flags().StringP("email", "e", "marco@rossi.com", "email of the author")
	err := signupCmd.MarkFlagRequired("email")
	cobra.CheckErr(err)

	signupCmd.Flags().StringP("password", "p", "password", "password of the author")
	err = signupCmd.MarkFlagRequired("password")
	cobra.CheckErr(err)

	rootCmd.AddCommand(signupCmd)
}
