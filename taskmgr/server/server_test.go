package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MrShanks/Taska/common/logger"
)

func TestGetTasksHandler(t *testing.T) {
	logger.InitLogger()
	t.Run("returns all tasks", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		taskHandler := TasksHandler{
			store: IMD,
		}

		taskHandler.ServeHTTP(response, request)

		got := response.Body.String()
		wantBytes, _ := json.Marshal(IMD.GetTasks())
		want := string(wantBytes)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}
