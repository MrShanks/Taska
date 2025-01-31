package server

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/MrShanks/Taska/taskmgr/logger"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	logger.InfoLogger.Println("Got request on / endpoint")
	io.WriteString(w, "Welcome to your dashboard\n")
}

func Listen() {
	webMux := http.NewServeMux()
	webMux.HandleFunc("/", HomeHandler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: webMux,
	}

	err := httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		logger.ErrorLogger.Println("Server closed")
	} else if err != nil {
		logger.ErrorLogger.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
