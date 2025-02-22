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
		newDummyTask := struct {
			Title string
			Desc  string
		}{
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
