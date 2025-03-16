package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/MrShanks/Taska/common/task"
)

func TestFetch(t *testing.T) {
	t.Run("Fetch Tasks returns all tasks in the store", func(t *testing.T) {
		tasks := []*task.Task{task.New("task1", "this is the task desc")}

		jsonTask, err := json.Marshal(tasks)
		if err != nil {
			t.Errorf("couldn't Marshal task")
		}

		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(jsonTask)
			if err != nil {
				t.Errorf("couldn't write the response")
			}
		}))
		defer mockServer.Close()

		serverURL, err := url.Parse(mockServer.URL)
		if err != nil {
			t.Errorf("couldn't parse mock server url")
		}

		mockClient := &Taskcli{
			HttpClient: &http.Client{},
			ServerURL:  *serverURL,
		}

		got := FetchTasks(mockClient, context.Background(), "/tasks", "token")
		want := jsonTask

		if reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
