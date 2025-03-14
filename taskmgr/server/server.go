package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/MrShanks/Taska/common/author"
	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/taskmgr/storage"
	"github.com/MrShanks/Taska/utils"
	"github.com/jackc/pgx/v5"
)

func NewServer(cfg *utils.Config, taskStore task.Store, authorStore author.Store) *http.Server {
	return &http.Server{
		Addr:              net.JoinHostPort(cfg.Spec.Host, cfg.Spec.Port),
		Handler:           InitMuxWithRoutes(taskStore, authorStore),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func ConnectDB(db_url string) (*pgx.Conn, error) {
	password := os.Getenv("POSTGRES_PWD")
	dburl := fmt.Sprintf(db_url, password)

	var err error

	conn, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Listen initialize the server and waits for requests
func Listen(cfg *utils.Config) error {
	var err error

	conn, err := ConnectDB(cfg.Spec.DB_URL)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	taskStore := storage.TaskStoreDB{Conn: conn}
	authorStore := storage.AuthorStoreDB{Conn: conn}

	httpServer := NewServer(cfg, &taskStore, &authorStore)
	log.Printf("Server version: %s listening at %s", cfg.Version, httpServer.Addr)

	err = httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	} else if !errors.Is(err, nil) {
		return err
	}
	return nil
}
