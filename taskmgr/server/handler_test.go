package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/MrShanks/Taska/common/task"
	"github.com/MrShanks/Taska/taskmgr/storage"
)

func TestNewTaskHandler(t *testing.T) {
	// Arrange
	mockDatabase := storage.InMemoryDatabase{
		Tasks: map[uuid.UUID]*task.Task{},
	}

	newDummyTask := task.New("new upcoming task", "title of a new fancy task")

	body, err := json.Marshal(newDummyTask)
	if err != nil {
		t.Errorf("Unable to marshal the dummy task")
	}

	// Act
	handler := NewTaskHandler(&mockDatabase)

	request, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/new", bytes.NewBuffer(body))
	response := httptest.NewRecorder()

	handler(response, request)

	// "Check if the pushed task has been created": mockDatabase.GetTasks(),
	got := mockDatabase.GetTasks()
	want := task.New(newDummyTask.Title, newDummyTask.Desc)

	if got[0].Title != want.Title {
		t.Errorf("got: %v, want: %v", got[0].Title, want.Title)
	}
}

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

func TestUpdateTaskHandler(t *testing.T) {
	t.Run("PUT request on /update/{task_id} updates the selected task", func(t *testing.T) {
		// Arrange
		task1 := task.New("before change", "before the changes have been applied")

		tasks := map[uuid.UUID]*task.Task{
			task1.ID: task1,
		}

		IMD := storage.InMemoryDatabase{
			Tasks: tasks,
		}

		// Act
		handler := UpdateTaskHandler(&IMD)

		body := &bytes.Buffer{}
		body.Write([]byte(`{"title":"new title","desc":"new description"}`))

		request, err := http.NewRequestWithContext(context.Background(), http.MethodPut, fmt.Sprintf("/delete/%s", task1.ID), body)
		if err != nil {
			t.Errorf("couldn't create request")
		}
		response := httptest.NewRecorder()

		handler(response, request)

		// Assert
		got := response.Body.String()
		want := fmt.Sprintf(`{"id":"%s","title":"new title","desc":"new description"}`, task1.ID.String())

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
			defer response.Result().Body.Close()

			test.handler(response, request)

			got := response.Code

			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
