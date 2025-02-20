package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
			t.Errorf("unable to fetch the tasks, task didn't created")
		}
		fmt.Printf("GOT  	%v\n", got[0].Title)
		fmt.Printf("WANT 	%v\n", want.Title)
		if got[0].Title != want.Title {
			t.Errorf("got: %v, want: %v", got[0].Title, want.Title)
		}
	})
}

// func TestNewTaskHandler(t *testing.T) {

// 	var IMD = inMemoryDatabase{
// 		tasks: []*task.Task{},
// 	}

// 	newDummyTask := struct {
// 		Title string
// 		Desc  string
// 	}{
// 		Title: "new upcomig task",
// 		Desc:  "title of a new fancy task",
// 	}

// 	t.Run("check new task created", func(t *testing.T) {

// 		jsonTask, err := json.Marshal(newDummyTask)
// 		if err != nil {
// 			t.Errorf("Unable to marshal the dummy task")
// 			return
// 		}

// 		request, _ := http.NewRequest(http.MethodPost, "/new", bytes.NewBuffer(jsonTask))

// 		body, err := io.ReadAll(request.Body)
// 		if err != nil {
// 			t.Errorf("Couldn't read the body. Error type: %s", err)
// 		}
// 		defer request.Body.Close()

// 		newTask := task.Task{
// 			ID: uuid.New(),
// 		}

// 		err = json.Unmarshal(body, &newTask)
// 		if err != nil {
// 			t.Errorf("Invalid JSON format")
// 			return
// 		}
// 		IMD.New(&newTask)

// 		if len(IMD.tasks) == 0 {
// 			t.Errorf("Expected a task to be added, but the task list is empty")
// 		} else {
// 			got := IMD.tasks[0].Title
// 			want := newDummyTask.Title

// 			if got != want {
// 				t.Errorf("got %q, want %q", got, want)
// 			}
// 		}

// 	})
// }
