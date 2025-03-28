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
	"github.com/MrShanks/Taska/utils"
)

const contentType = "Content-Type"
const appJson = "application/json"

func GetOneTaskHandler(taskStore task.Store, authorStore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /task endpoint")

		taskID := r.PathValue("id")

		authorID := tokenToAuthor(r, authorStore)

		selectedTask, err := taskStore.GetOne(taskID, authorID)
		if err != nil {
			log.Printf("Couldn't retrieve task from store: %v", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(nil))
			return
		}

		jsonTask, err := json.Marshal(selectedTask)
		if err != nil {
			log.Printf("Couldn't Marshal tasks into json format: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(nil))
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

func GetAllTasksHandler(taskStore task.Store, authorStore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /tasks endpoint")

		authorID := tokenToAuthor(r, authorStore)

		tasks := taskStore.GetTasks(authorID)

		jsonTasks, err := json.Marshal(tasks)
		if err != nil {
			log.Printf("Couldn't Marshal tasks into json format: %v", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(nil))
			return
		}

		w.Header().Set(contentType, appJson)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonTasks)
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}
	}
}

func NewTaskHandler(taskStore task.Store, authorStore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /new endpoint")

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

		authorID := tokenToAuthor(r, authorStore)

		newTask.AuthorID = uuid.MustParse(authorID)

		newTaskID := taskStore.New(&newTask)
		if newTaskID == uuid.Nil {
			log.Printf("error, task with same name already exists")
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("error, task with same name already exists"))
			if err != nil {
				log.Printf("Couldn't write response: %v", err)
			}
			return
		}

		log.Printf("New task created with ID: %s", newTaskID)

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(newTaskID.String()))
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}

		EventLogger.WriteNew(newTask.ID, newTask.Title, newTask.Desc)
	}
}

func UpdateTaskHandler(taskStore task.Store, authorStore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /update endpoint")

		taskID := r.PathValue("id")
		authorID := tokenToAuthor(r, authorStore)

		changes := task.Task{}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Couldn't read the body. Error type: %s", err)
		}
		defer r.Body.Close()

		err = json.Unmarshal(reqBody, &changes)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(nil))
			log.Printf("Couldn't unmarshal the body. Error type: %s", err)
			return
		}

		updatedTask, err := taskStore.Update(taskID, changes.Title, changes.Desc, authorID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(nil))
			log.Printf("Error updating task: %v", err)
			return
		}

		jsonTask, err := json.Marshal(updatedTask)
		if err != nil {
			log.Printf("Couldn't marshal task into json format: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(nil))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonTask)
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}

		EventLogger.WriteMod(updatedTask.ID, updatedTask.Title, updatedTask.Desc)
	}
}

func DeleteTaskHandler(taskStore task.Store, authorStore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /delete endpoint")

		taskID := r.PathValue("id")

		authorID := tokenToAuthor(r, authorStore)

		err := taskStore.Delete(taskID, authorID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Printf("Deletion failed: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(nil))

		EventLogger.WriteDel(uuid.MustParse(taskID))
	}
}

func SearchTaskHandler(taskStore task.Store, authorStore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /search endpoint")

		keyword := r.PathValue("keyword")

		authorID := tokenToAuthor(r, authorStore)

		tasks, err := taskStore.Search(keyword, authorID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(nil))
			log.Printf("Couldn't find any task: %v", err)
			return
		}

		jsonTasks, err := json.Marshal(tasks)
		if err != nil {
			log.Printf("Couldn't Marshal tasks into json format: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(nil))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonTasks)
		if err != nil {
			log.Printf("Couldn't write response: %v", err)
		}
	}
}

func ImportTaskHandler(taskStore task.Store, authorStore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /import endpoint")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(nil))
			return
		}

		var tasks []*task.Task

		cType := r.Header.Get(contentType)
		if cType == appJson {
			err = json.Unmarshal(body, &tasks)
			if err != nil {
				log.Printf("Couldn't unmarshal tasks into a slice of tasks: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(nil))
				return
			}
		} else if cType == "application/x-yaml" {
			err = yaml.Unmarshal(body, &tasks)
			if err != nil {
				log.Printf("Couldn't unmarshal tasks into a slice of tasks: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(nil))
				return
			}
		}

		authorID := tokenToAuthor(r, authorStore)

		taskStore.BulkImport(tasks, authorID)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(nil))
	}
}

func Signup(authorStore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /signup endpoint")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(nil))
			log.Printf("Couldn't read the body. Error type: %v", err)
			return
		}
		defer r.Body.Close()

		newAuthor := author.Author{}

		err = json.Unmarshal(body, &newAuthor)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(nil))
			log.Printf("Couldn't unmarshal the body. Error type: %s", err)
			return
		}

		if err = authorStore.SignUp(&newAuthor); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(nil))
			log.Printf("Couldn't sign you up: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte("Signup successful!")); err != nil {
			log.Printf("Couldn't write signup response: %v", err)
		}
	}
}

func Signin(authorStore author.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request on /signin endpoint")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(nil))
			log.Printf("Couldn't read the body. Error type: %v", err)
			return
		}
		defer r.Body.Close()

		signInAuthor := author.Author{}

		err = json.Unmarshal(body, &signInAuthor)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(nil))
			log.Printf("Couldn't unmarshal the body. Error type: %s", err)
			return
		}

		token, err := utils.CreateToken(signInAuthor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(nil))
			log.Printf("Error creating token: %v", err)
			return
		}

		if err = authorStore.SignIn(signInAuthor.Email, signInAuthor.Password, token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(nil))
			log.Printf("Error during authentication: %v", err)
			return
		}

		LoggedAuthorToken = token

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
	w.Write([]byte(nil))
}

func tokenToAuthor(r *http.Request, authorStore author.Store) string {
	token := r.Header.Get("Token")

	authorID, err := authorStore.GetAuthorID(token)
	if err != nil {
		log.Printf("Error when fetching the author id: %v", err)
	}

	return authorID
}
