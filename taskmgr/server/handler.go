package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/MrShanks/Taska/common/task"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

func GetOneTaskHandler(store task.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Printf("Method: %s is not allowed on /task endpoint\n", r.Method)

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Printf("Got request on /task endpoint\n")

		taskID := strings.Split(r.URL.Path, "/")[2]

		selectedTask, err := store.GetOne(taskID)
		if err != nil {
			log.Printf("Couldn't retrieve task from store: %v\n", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jsonTask, err := json.Marshal(selectedTask)
		if err != nil {
			log.Printf("Couldn't Marshal tasks into json format: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonTask)
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}
	}
}

func GetAllTasksHandler(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := isAllowedMethod(http.MethodGet, w, r); err != nil {
			return
		}

		log.Printf("Got request on /tasks endpoint\n")

		tasks := store.GetTasks()
		if tasks == nil {
			_, err := w.Write([]byte("Could't able to reach database"))
			if err != nil {
				log.Printf("Error: %v", err)
			}
			return
		}
		jsonTasks, err := json.Marshal(tasks)
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
		if err := isAllowedMethod(http.MethodPost, w, r); err != nil {
			return
		}

		log.Printf("Got request on /new endpoint\n")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Couldn't read the body. Error type: %s", err)
		}
		defer r.Body.Close()

		newTask := task.Task{}

		err = json.Unmarshal(body, &newTask)
		if err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
		newTaskID := store.New(&newTask)
		if newTaskID == uuid.Nil {
			_, err := w.Write([]byte("Could't able to reach database"))
			if err != nil {
				log.Printf("Error: %v", err)
			}
			return
		}

		log.Printf("New task created. ID: %s", newTaskID)

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(newTaskID.String()))
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}
	}
}

func UpdateTaskHandler(store task.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := isAllowedMethod(http.MethodPut, w, r); err != nil {
			return
		}

		log.Printf("Got request on /update endpoint\n")

		taskID := strings.Split(r.URL.Path, "/")[2]

		changes := task.Task{}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Couldn't read the body. Error type: %s", err)
		}
		defer r.Body.Close()

		err = json.Unmarshal(reqBody, &changes)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldn't unmarshal the body. Error type: %s", err)
			return
		}

		updatedTask, err := store.Update(taskID, changes.Title, changes.Desc)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("Error updating task: %v", err)
			return
		}

		jsonTask, err := json.Marshal(updatedTask)
		if err != nil {
			log.Printf("Couldn't marshal task into json format: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write(jsonTask)
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}
	}
}

func DeleteTaskHandler(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := isAllowedMethod(http.MethodDelete, w, r); err != nil {
			return
		}

		log.Printf("Got request on /delete endpoint\n")

		taskID := strings.Split(r.URL.Path, "/")[2]

		err := store.Delete(taskID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("Deletion failed: %v", err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func ImportTaskHandler(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := isAllowedMethod(http.MethodPost, w, r); err != nil {
			return
		}

		log.Printf("Got request on /import endpoint\n")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var tasks []*task.Task

		contentType := r.Header.Get("Content-Type")
		if contentType == "application/json" {
			err = json.Unmarshal(body, &tasks)
			if err != nil {
				log.Printf("Couldn't unmarshal tasks into a slice of tasks: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else if contentType == "application/x-yaml" {
			err = yaml.Unmarshal(body, &tasks)
			if err != nil {
				log.Printf("Couldn't unmarshal tasks into a slice of tasks: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		store.BulkImport(tasks)

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

func isAllowedMethod(httpMethod string, w http.ResponseWriter, r *http.Request) error {
	if r.Method != httpMethod {
		log.Printf("Method: %s is not allowed\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return nil
}
