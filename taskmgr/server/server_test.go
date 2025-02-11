package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MrShanks/Taska/common/logger"
	"github.com/MrShanks/Taska/common/task"
)

func TestGetTasksHandler(t *testing.T) {
	logger.InitLogger()

	var IMD = inMemoryDatabase{
		tasks: []*task.Task{
			task.New("first", "Desc First"),
			task.New("second", "Desc Second"),
			task.New("third", "Desc Third"),
		},
	}

	t.Run("returns all tasks", func(t *testing.T) {
		handler := GetHandler(&IMD)

		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		handler(response, request)

		got := response.Body.String()
		wantBytes, _ := json.Marshal(IMD.GetTasks())
		want := string(wantBytes)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}
