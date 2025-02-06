package task

import (
	"github.com/google/uuid"
)

type Task struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Desc  string    `json:"desc"`
}

func New(title, desc string) *Task {
	return &Task{
		ID:    uuid.New(),
		Title: title,
		Desc:  desc,
	}
}

type TaskStore interface {
	GetTasks() []*Task
}
