package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/MrShanks/Taska/common/logger"
	"github.com/MrShanks/Taska/common/task"
)

var IMD = inMemoryDatabase{
	tasks: []*task.Task{
		task.New("first", "Desc First"),
		task.New("second", "Desc Second"),
		task.New("third", "Desc Third"),
	},
}

// inMemoryDatabase implements the taskStore interface
type inMemoryDatabase struct {
	tasks []*task.Task
}

func (md inMemoryDatabase) GetTasks() []*task.Task {
	return md.tasks
}

// TasksHandler implements the Handler interface it takes a store of type task.TaskStore
type TasksHandler struct {
	store task.TaskStore
}

func (t TasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.ErrorLogger.Printf("Method: %s is not allowed on /tasks endpoint\n", r.Method)

		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	logger.InfoLogger.Printf("Got request on /tasks endpoint\n")

	jsonTasks, err := json.Marshal(t.store.GetTasks())
	if err != nil {
		logger.ErrorLogger.Printf("Couldn't Marshal tasks into json format: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonTasks)
	if err != nil {
		logger.ErrorLogger.Printf("Couldn't write response: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("Got request on / endpoint")
	_, err := io.WriteString(w, "Welcome to your dashboard\n")
	if err != nil {
		logger.ErrorLogger.Printf("Couldn't write response: %v", err)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	// Just return a blank 204 response (No Content)
	// This is needed because browser always issue a second
	// http request to get the favicon for a website
	w.WriteHeader(http.StatusNoContent)
}

// Listen initialize the server and waits for requests
func Listen(version string) {
	tasksHandler := TasksHandler{store: IMD}

	webMux := http.NewServeMux()
	webMux.HandleFunc("/", homeHandler)
	webMux.Handle("/tasks", tasksHandler)
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
