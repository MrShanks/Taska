package tskmgt

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got request on /")
	io.WriteString(w, "Welcome to your dashboard")
}

func Listner() {
	http.HandleFunc("/", Root)

	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
