package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/MrShanks/Taska/common/logger"
	"github.com/MrShanks/Taska/common/model"
)

// In memory database of tasks until we have a proper database connected
var tasks = []*model.Task{
	model.New("First", "First description"),
	model.New("Second", "Second description"),
	model.New("Third", "Third description"),
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("Got request on / endpoint")
	io.WriteString(w, "Welcome to your dashboard\n")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	// Just return a blank 204 response (No Content)
	// This is needed because browser always issue a second
	// http request to get the favicon for a website
	w.WriteHeader(http.StatusNoContent)
}

// ListTaskHandler returns all tasks
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.ErrorLogger.Printf("Method: %s is not allowed on /tasks endpoint\n", r.Method)

		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	logger.InfoLogger.Printf("Got request on /tasks endpoint\n")

	jsonTasks, err := json.Marshal(tasks)
	if err != nil {
		logger.ErrorLogger.Printf("Couldn't Marshal tasks into json format: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonTasks)
}

// Listen initialize the server and waits for requests
func Listen(version string) {
	webMux := http.NewServeMux()
	webMux.HandleFunc("/", HomeHandler)
	webMux.HandleFunc("/tasks", GetTasksHandler)
	webMux.HandleFunc("/favicon.ico", faviconHandler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: webMux,
	}

	logger.InfoLogger.Printf("Server version: %s listening at %s", version, httpServer.Addr)

	err := httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		logger.ErrorLogger.Println("Server closed")
	} else if err != nil {
		logger.ErrorLogger.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
