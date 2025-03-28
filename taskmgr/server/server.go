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
	"github.com/MrShanks/Taska/taskmgr/logger"
	"github.com/MrShanks/Taska/taskmgr/storage"
	"github.com/MrShanks/Taska/utils"
	"github.com/jackc/pgx/v5"
)

var EventLogger logger.TransactionLogger
var LoggedAuthorToken string

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
			log.Printf("%v", err)
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

func ConnectDB(ctx context.Context, dbURL string) (*pgx.Conn, error) {
	password := os.Getenv("POSTGRES_PWD")
	user := os.Getenv("DBUSER")
	database := os.Getenv("DATABASE")
	host := os.Getenv("HOST")
	dburl := fmt.Sprintf(dbURL, user, password, host, database)

	conn, err := pgx.Connect(ctx, dburl)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Listen initializes the server the storage and the transaction log and then waits for requests
func Listen(cfg *utils.Config) error {
	var err error

	ctxConnection, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := ConnectDB(ctxConnection, cfg.Spec.DB_URL)
	if err != nil {
		return err
	}
	defer conn.Close(ctxConnection)

	taskStore := storage.TaskStore{Conn: conn}
	authorStore := storage.AuthorStore{Conn: conn}

	httpServer := NewServer(cfg, &taskStore, &authorStore)

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
