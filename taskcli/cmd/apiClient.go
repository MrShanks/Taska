package cmd

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/MrShanks/Taska/utils"
)

type Tasckli struct {
	HttpClient *http.Client
	Cfg        *utils.Config
	ServerURL  url.URL
}

func NewApiClient() *Tasckli {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	cfg := utils.LoadConfig("config.yaml")

	serverURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(cfg.Spec.Host, cfg.Spec.Port),
	}

	return &Tasckli{
		HttpClient: httpClient,
		Cfg:        cfg,
		ServerURL:  serverURL,
	}
}
