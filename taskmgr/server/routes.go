package server

import (
	"net/http"

	"github.com/MrShanks/Taska/common/author"
	"github.com/MrShanks/Taska/common/task"
)

func InitMuxWithRoutes(taskStore task.Store, authorStore author.Store) *http.ServeMux {
	webMux := http.NewServeMux()
	webMux.HandleFunc("GET /", homeHandler)
	webMux.HandleFunc("GET /favicon.ico", faviconHandler)

	// Tasks related Routes
	webMux.HandleFunc("GET /tasks", LoggedInOnly(GetAllTasksHandler(taskStore, authorStore)))
	webMux.HandleFunc("GET /task/{id}", LoggedInOnly(GetOneTaskHandler(taskStore, authorStore)))
	webMux.HandleFunc("GET /search/{keyword}", LoggedInOnly(SearchTaskHandler(taskStore, authorStore)))
	webMux.HandleFunc("POST /new", LoggedInOnly(NewTaskHandler(taskStore, authorStore)))
	webMux.HandleFunc("POST /import", LoggedInOnly(ImportTaskHandler(taskStore, authorStore)))
	webMux.HandleFunc("DELETE /delete/{id}", LoggedInOnly(DeleteTaskHandler(taskStore, authorStore)))
	webMux.HandleFunc("PUT /update/{id}", LoggedInOnly(UpdateTaskHandler(taskStore, authorStore)))

	// Users related Routes
	webMux.HandleFunc("POST /signup", Signup(authorStore))
	webMux.HandleFunc("POST /signin", Signin(authorStore))

	return webMux
}
