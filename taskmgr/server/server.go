package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/taskmgr/logger"
	"github.com/MrShanks/Taska/taskmgr/storage"
	"github.com/MrShanks/Taska/utils"
)

var EventLogger logger.TransactionLogger

func initTransactionLog() error {
	var err error

	EventLogger, err = logger.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("couldn't create event logger: %v", err)
	}

	EventLogger.Run()

	return err
}

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

// Listen initializes the server the storage and the transaction log and then waits for requests
func Listen(cfg *utils.Config) {
	IMD := storage.InMemoryDatabase{
		Tasks: []*task.Task{
			task.New("first", "Desc First"),
			task.New("second", "Desc Second"),
			task.New("third", "Desc Third"),
		},
	}

	err := initTransactionLog()
	if err != nil {
		log.Printf("error occured during transaction log initialization: %v", err)
		os.Exit(1)
	}

	httpServer := NewServer(cfg, &IMD)

	log.Printf("Server version: %s listening at %s", cfg.Version, httpServer.Addr)

	err = httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if !errors.Is(err, nil) {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
