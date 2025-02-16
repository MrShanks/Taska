package cmd

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/MrShanks/Taska/utils"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get all active tasks",
	Long:  "get all active tasks store on the servere",

	Run: func(cmd *cobra.Command, args []string) {
		apiClient := &http.Client{}
		cfg := utils.LoadConfig("config.yaml")

		serverURL := url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(cfg.Spec.Host, cfg.Spec.Port),
			Path:   "/tasks",
		}

		response, err := apiClient.Get(serverURL.String())
		if err != nil {
			log.Printf("Couldn't get a response from the server: %v", err)
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("Couldn't read response body: %v", err)
		}

		fmt.Printf("%v\n", string(bodyBytes))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
