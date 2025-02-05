package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MrShanks/Taska/common/logger"
)

func TestGetTasksHanlder(t *testing.T) {
	logger.InitLogger()
	t.Run("returns all tasks", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		GetTasksHandler(response, request)

		got := response.Body.String()
		wantBytes, _ := json.Marshal(tasks)
		want := string(wantBytes)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}
