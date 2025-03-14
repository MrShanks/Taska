package server

import (
	"net/http"

	"github.com/MrShanks/Taska/common/task"
)

func InitMuxWithRoutes(store task.Store) *http.ServeMux {
	webMux := http.NewServeMux()
	webMux.HandleFunc("/", homeHandler)
	webMux.HandleFunc("/favicon.ico", faviconHandler)

	// Tasks related Routes
	webMux.HandleFunc("/tasks", GetAllTasksHandler(store))
	webMux.HandleFunc("/task/{task_id}", GetOneTaskHandler(store))
	webMux.HandleFunc("/new", NewTaskHandler(store))
	webMux.HandleFunc("/import", ImportTaskHandler(store))
	webMux.HandleFunc("/delete/{task_id}", DeleteTaskHandler(store))
	webMux.HandleFunc("/update/{task_id}", UpdateTaskHandler(store))

	return webMux
}
