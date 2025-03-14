package task

import (
	"github.com/MrShanks/Taska/common/task"
	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Tasks     []*task.Task
}

func New(first, last, email, password string) *User {
	return &User{
		UUID:      uuid.New(),
		Firstname: first,
		Lastname:  last,
		Email:     email,
		Password:  password,
		Tasks:     []*task.Task{},
	}
}
