package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/taskmgr/storage"
)

func TestNewTaskHandler(t *testing.T) {
	t.Run("Create new task and check if it's there", func(t *testing.T) {

		// Arrange
		mockDatabase := storage.InMemoryDatabase{
			Tasks: []*task.Task{},
		}
		newDummyTask := task.Task{
			Title: "new upcomig task",
			Desc:  "title of a new fancy task",
		}

		body, err := json.Marshal(newDummyTask)
		if err != nil {
			t.Errorf("Unable to marshal the dummy task")
			return
		}

		// Act
		handler := NewTaskHandler(&mockDatabase)

		request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/new", bytes.NewBuffer(body))
		response := httptest.NewRecorder()

		handler(response, request)

		// Assert
		want := task.New(newDummyTask.Title, newDummyTask.Desc)
		got := mockDatabase.GetTasks()
		if len(got) == 0 {
			t.Errorf("Expected got != %v", len(got))
			return
		}
		t.Logf("GOT  	%v\n", got[0].Title)
		t.Logf("WANT 	%v\n", want.Title)
		if got[0].Title != want.Title {
			t.Errorf("got: %v, want: %v", got[0].Title, want.Title)
			return
		}
	})
}

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
