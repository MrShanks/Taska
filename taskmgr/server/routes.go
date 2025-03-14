package server

import (
	"net/http"

	"github.com/MrShanks/Taska/common/author"
	"github.com/MrShanks/Taska/common/task"
)

func InitMuxWithRoutes(taskStore task.Store, authorStore author.Store) *http.ServeMux {
	webMux := http.NewServeMux()
	webMux.HandleFunc("/", homeHandler)
	webMux.HandleFunc("/favicon.ico", faviconHandler)

	// Tasks related Routes
	webMux.HandleFunc("/tasks", GetAllTasksHandler(taskStore))
	webMux.HandleFunc("/task/{task_id}", GetOneTaskHandler(taskStore))
	webMux.HandleFunc("/new", NewTaskHandler(taskStore))
	webMux.HandleFunc("/import", ImportTaskHandler(taskStore))
	webMux.HandleFunc("/delete/{task_id}", DeleteTaskHandler(taskStore))
	webMux.HandleFunc("/update/{task_id}", UpdateTaskHandler(taskStore))

	// Users related Routes
	webMux.HandleFunc("/signup", Signup(authorStore))
	webMux.HandleFunc("/signin", Signin(authorStore))

	return webMux
}
