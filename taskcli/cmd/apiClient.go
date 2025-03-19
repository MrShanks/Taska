package cmd

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/MrShanks/Taska/utils"
)

type Taskcli struct {
	HttpClient *http.Client
	Cfg        *utils.Config
	ServerURL  url.URL
	Token      string
}

func NewApiClient() *Taskcli {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	cfg := utils.LoadConfig("config.yaml")

	serverURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(cfg.Spec.Host, cfg.Spec.Port),
	}

	return &Taskcli{
		HttpClient: httpClient,
		Cfg:        cfg,
		ServerURL:  serverURL,
	}
}
