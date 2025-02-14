package server

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/taskmgr/storage"
)

// Listen initialize the server and waits for requests
func Listen(version string) {
	IMD := storage.InMemoryDatabase{
		Tasks: []*task.Task{
			task.New("first", "Desc First"),
			task.New("second", "Desc Second"),
			task.New("third", "Desc Third"),
		},
	}

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           InitMuxWithRoutes(&IMD),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Printf("Server version: %s listening at %s", version, httpServer.Addr)

	err := httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if !errors.Is(err, nil) {
		log.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
