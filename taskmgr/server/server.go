package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/MrShanks/Taska/common/logger"
	"github.com/MrShanks/Taska/common/task"
	"github.com/google/uuid"
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

func (md *inMemoryDatabase) GetTasks() []*task.Task {
	return md.tasks
}

func (md *inMemoryDatabase) New(task *task.Task) {
	md.tasks = append(md.tasks, task)
}

func GetHandler(store task.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			logger.ErrorLogger.Printf("Method: %s is not allowed on /tasks endpoint\n", r.Method)

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		logger.InfoLogger.Printf("Got request on /tasks endpoint\n")

		jsonTasks, err := json.Marshal(store.GetTasks())
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
}

func NewTaskHandler(store task.TaskStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.ErrorLogger.Printf("Method: %s is not allowed on /new endpoint\n", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		logger.InfoLogger.Printf("Got request on /new endpoint\n")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.ErrorLogger.Printf("Couldn't read the body. Error type: %s", err)
		}
		defer r.Body.Close()

		newTask := task.Task{
			ID: uuid.New(),
		}

		err = json.Unmarshal(body, &newTask)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		store.New(&newTask)

		logger.InfoLogger.Printf("New task created.\nID: %s", newTask.ID)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("New task created.\nID: %s\nTitle: %s\nDesc: %s", newTask.ID, newTask.Title, newTask.Desc)))
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
	webMux := http.NewServeMux()
	webMux.HandleFunc("/", homeHandler)
	webMux.HandleFunc("/tasks", GetHandler(&IMD))
	webMux.HandleFunc("/new", NewTaskHandler(&IMD))
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
