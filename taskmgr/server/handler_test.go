package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/taskmgr/storage"
)

func TestGetTasksHandler(t *testing.T) {
	t.Run("GET request on /tasks returns all tasks", func(t *testing.T) {
		// Arrange
		task1 := task.New("first", "Desc First")
		task2 := task.New("second", "Desc Second")

		tasks := map[uuid.UUID]*task.Task{
			task1.ID: task1,
			task2.ID: task2,
		}

		IMD := storage.InMemoryDatabase{
			Tasks: tasks,
		}

		// Act
		handler := GetAllTasksHandler(&IMD)

		request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks", nil)
		if err != nil {
			t.Errorf("couldn't create request")
		}
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

			test.handler(response, request)

			got := response.Result().StatusCode

			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
