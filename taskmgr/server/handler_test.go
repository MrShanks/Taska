package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeneralHandler(t *testing.T) {

	tests := []struct {
		desc     string
		want     int
		endpoint string
		handler  http.HandlerFunc
	}{
		{"GET request on / returns status 200", 200, "/", homeHandler},
		{"GET request on /favicon.ico returns status 204", 204, "/favicon.ico", faviconHandler},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			request, err := http.NewRequestWithContext(context.Background(), "GET", test.endpoint, nil)
			if err != nil {
				t.Errorf("couldn't create request")
			}

			response := httptest.NewRecorder()
			defer response.Result().Body.Close()

			test.handler(response, request)

			got := response.Code

			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
