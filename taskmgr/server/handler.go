package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/MrShanks/Taska/common/task"
	"github.com/google/uuid"
)

func GetAllTasksHandler(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Printf("Method: %s is not allowed on /tasks endpoint\n", r.Method)

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Printf("Got request on /tasks endpoint\n")

		jsonTasks, err := json.Marshal(store.GetTasks())
		if err != nil {
			log.Printf("Couldn't Marshal tasks into json format: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonTasks)
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}
	}
}

func NewTaskHandler(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Printf("Method: %s is not allowed on /new endpoint\n", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Printf("Got request on /new endpoint\n")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Couldn't read the body. Error type: %s", err)
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

		log.Printf("New task created. ID: %s", newTask.ID)
		w.WriteHeader(http.StatusCreated)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request on / endpoint")

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Welcome to your dashboard"))
	if err != nil {
		log.Printf("Couldn't write response: %v", err)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	// Returns a blank 204 response (No Content)
	// This is needed because browser always issue a second
	// http request to get the favicon for a website
	w.WriteHeader(http.StatusNoContent)
}
