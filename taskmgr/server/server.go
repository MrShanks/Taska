package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
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

	events, errs := EventLogger.ReadEvents()

	e, ok := logger.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errs:
		case e, ok = <-events:
			switch e.Type {
			case logger.Del:
				// TODO: implement Deletion
			case logger.Mod:
				// TODO: implement Modification
			case logger.New:
				// TODO: implement Creation
			}
		}
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
func Listen(cfg *utils.Config) error {
	var err error
	DB := storage.PostgresDatabase{}

	err = DB.Connect(cfg.Spec.DB_URL)
	if err != nil {
		return err
	}
	defer DB.Conn.Close(context.Background())

	httpServer := NewServer(cfg, &DB)

	err = initTransactionLog()
	if err != nil {
		log.Printf("error occured during transaction log initialization: %v", err)
		return err
	}

	log.Printf("Server version: %s listening at %s", cfg.Version, httpServer.Addr)

	err = httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}
