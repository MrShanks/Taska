package task

import "github.com/google/uuid"

type Store interface {
	GetTasks() map[uuid.UUID]*Task
	New(*Task) uuid.UUID
	Delete(id string) error
	Update(id, title, desc string) (*Task, error)
}
