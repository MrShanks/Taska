package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/taskmgr/storage"
)

func TestGetTasksHandler(t *testing.T) {
	t.Run("GET request on /tasks returns all tasks", func(t *testing.T) {

		// Arrange
		IMD := storage.InMemoryDatabase{
			Tasks: []*task.Task{
				task.New("first", "Desc First"),
				task.New("second", "Desc Second"),
				task.New("third", "Desc Third"),
			},
		}

		// Act
		handler := GetAllTasksHandler(&IMD)

		request, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		handler(response, request)

		// Assert
		got := response.Body.String()
		wantBytes, _ := json.Marshal(IMD.GetTasks())
		want := string(wantBytes)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
