package server

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

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
func Listen(cfg *utils.Config) {

	task1 := task.New("first", "Desc First")
	task2 := task.New("second", "Desc Second")
	task3 := task.New("third", "Desc Third")

	tasks := map[uuid.UUID]*task.Task{
		task1.ID: task1,
		task2.ID: task2,
		task3.ID: task3,
	}

	IMD := storage.InMemoryDatabase{
		Tasks: tasks,
	}

	httpServer := NewServer(cfg, &IMD)

	log.Printf("Server version: %s listening at %s", cfg.Version, httpServer.Addr)

	err := httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if !errors.Is(err, nil) {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
