package server

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
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

func TestNewTaskHandler(t *testing.T) {

	var IMD = inMemoryDatabase{
		tasks: []*task.Task{},
	}

	newDummyTask := struct {
		Title string
		Desc  string
	}{
		Title: "new upcomig task",
		Desc:  "title of a new fancy task",
	}

	t.Run("check new task created", func(t *testing.T) {

		jsonTask, err := json.Marshal(newDummyTask)
		if err != nil {
			t.Errorf("Unable to marshal the dummy task")
			return
		}

		request, _ := http.NewRequest(http.MethodPost, "/new", bytes.NewBuffer(jsonTask))

		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Errorf("Couldn't read the body. Error type: %s", err)
		}
		defer request.Body.Close()

		newTask := task.Task{
			ID: uuid.New(),
		}

		err = json.Unmarshal(body, &newTask)
		if err != nil {
			t.Errorf("Invalid JSON format")
			return
		}
		IMD.New(&newTask)

		if len(IMD.tasks) == 0 {
			t.Errorf("Expected a task to be added, but the task list is empty")
		} else {
			got := IMD.tasks[0].Title
			want := newDummyTask.Title

			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		}

	})
}
