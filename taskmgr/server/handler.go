package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"

	"github.com/MrShanks/Taska/common/author"
	"github.com/MrShanks/Taska/common/task"
)

const contentType = "Content-Type"
const appJson = "application/json"

func GetOneTaskHandler(store task.Store) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /task endpoint\n")

		taskID := r.PathValue("id")

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

		w.Header().Set(contentType, appJson)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonTask)
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}
	}
}

func GetAllTasksHandler(store task.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /tasks endpoint\n")

		tasks := store.GetTasks()

		jsonTasks, err := json.Marshal(tasks)
		if err != nil {
			log.Printf("Couldn't Marshal tasks into json format: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set(contentType, appJson)
		_, err = w.Write(jsonTasks)
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}
	}
}

func NewTaskHandler(store task.Store, astore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		token := r.Header.Get("Token")

		authorID, err := astore.GetAuthorID(token)
		if err != nil {
			log.Printf("error fetching author id: %v", err)
		}

		newTask.AuthorID = uuid.MustParse(authorID)

		newTaskID := store.New(&newTask)
		if newTaskID == uuid.Nil {
			_, err := w.Write([]byte("Couldn't reach the database"))
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
		log.Printf("Got request on /update endpoint\n")

		taskID := r.PathValue("id")

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
		log.Printf("Got request on /delete endpoint\n")

		taskID := r.PathValue("id")

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
		log.Printf("Got request on /import endpoint\n")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var tasks []*task.Task

		cType := r.Header.Get(contentType)
		if cType == appJson {
			err = json.Unmarshal(body, &tasks)
			if err != nil {
				log.Printf("Couldn't unmarshal tasks into a slice of tasks: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else if cType == "application/x-yaml" {
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

func Signup(store author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /signup endpoint\n")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldn't read the body. Error type: %v", err)
			return
		}
		defer r.Body.Close()

		newAuthor := author.Author{}

		err = json.Unmarshal(body, &newAuthor)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldn't unmarshal the body. Error type: %s", err)
			return
		}

		if err = store.SignUp(&newAuthor); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldn't sign you up: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte("Signup successful!")); err != nil {
			log.Printf("Couldn't write signup response: %v", err)
		}
	}
}

func Signin(store author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /signin endpoint")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldn't read the body. Error type: %v", err)
			return
		}
		defer r.Body.Close()

		signInAuthor := author.Author{}

		err = json.Unmarshal(body, &signInAuthor)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Couldn't unmarshal the body. Error type: %s", err)
			return
		}

		token := uuid.New().String()
		if err = store.SignIn(signInAuthor.Email, signInAuthor.Password, token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("Error during authentication: %v", err)
			return
		}

		loggedAuthors[token] = signInAuthor.Email

		w.Header().Set("token", token)
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("Login successful!")); err != nil {
			log.Printf("Couldn't write login response: %v", err)
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

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
