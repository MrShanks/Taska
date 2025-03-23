package cmd

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

// Mock HTTP Client by implementing RoundTripper
type mockRoundTripper struct {
	mockResponse *http.Response
	mockError    error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockResponse, m.mockError
}

func TestDelTask(t *testing.T) {
	tests := []struct {
		name       string
		mockResp   *http.Response
		mockErr    error
		want       string
		statusCode int
	}{
		{
			name: "Delete task should return task successffully deleted",
			mockResp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       http.NoBody,
			},
			want: "Task Successfully deleted\n",
		},
		{
			name: "Delete task that doesn't exist returns task not found",
			mockResp: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       http.NoBody,
			},
			want: "Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &http.Client{
				Transport: &mockRoundTripper{
					mockResponse: tt.mockResp,
					mockError:    tt.mockErr,
				},
			}

			taskcli := &Taskcli{
				ServerURL:  url.URL{Scheme: "http", Host: "localhost"},
				HttpClient: mockClient,
			}

			got := delTask(taskcli, context.Background(), "/delete", "notoken")
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
