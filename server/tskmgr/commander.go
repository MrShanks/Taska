package tskmgr

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got request on / endpoint")
	io.WriteString(w, "Welcome to your dashboard")
}

func Listner() {
	webMux := http.NewServeMux()
	webMux.HandleFunc("/", Root)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: webMux,
	}

	err := httpServer.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
