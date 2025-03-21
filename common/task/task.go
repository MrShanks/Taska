package task

import "github.com/google/uuid"

type Task struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Desc     string    `json:"desc"`
	AuthorID uuid.UUID
}

func New(title, desc string) *Task {
	return &Task{
		Title: title,
		Desc:  desc,
	}
}
