package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/taskmgr/storage"
	"github.com/MrShanks/Taska/utils"
)

func NewServer(cfg *utils.Config, store task.Store) *http.Server {
	return &http.Server{
		Addr:              net.JoinHostPort(cfg.Spec.Host, cfg.Spec.Port),
		Handler:           InitMuxWithRoutes(store),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

// Listen initialize the server and waits for requests
func Listen(cfg *utils.Config) error {
	var err error
	DB := storage.PostgresDatabase{}

	err = DB.Connect(cfg.Spec.DB_URL)
	if err != nil {
		return err
	}
	defer DB.Conn.Close(context.Background())

	httpServer := NewServer(cfg, &DB)
	log.Printf("Server version: %s listening at %s", cfg.Version, httpServer.Addr)

	err = httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	} else if !errors.Is(err, nil) {
		return err
	}
	return nil
}
