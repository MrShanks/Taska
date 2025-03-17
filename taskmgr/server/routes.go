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
	webMux.HandleFunc("/tasks", LoggedInOnly(GetAllTasksHandler(taskStore)))
	webMux.HandleFunc("/task/{task_id}", LoggedInOnly(GetOneTaskHandler(taskStore)))
	webMux.HandleFunc("/new", LoggedInOnly(NewTaskHandler(taskStore)))
	webMux.HandleFunc("/import", LoggedInOnly(ImportTaskHandler(taskStore)))
	webMux.HandleFunc("/delete/{task_id}", LoggedInOnly(DeleteTaskHandler(taskStore)))
	webMux.HandleFunc("/update/{task_id}", LoggedInOnly(UpdateTaskHandler(taskStore)))

	// Users related Routes
	webMux.HandleFunc("/signup", Signup(authorStore))
	webMux.HandleFunc("/signin", Signin(authorStore))

	return webMux
}
