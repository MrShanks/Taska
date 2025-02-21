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
	webMux.HandleFunc("/new", NewTaskHandler(store))

	return webMux
}
