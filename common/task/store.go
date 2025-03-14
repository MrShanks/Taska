package task

import "github.com/google/uuid"

type Store interface {
	GetTasks() []*Task
	GetOne(id string) (*Task, error)
	New(*Task) uuid.UUID
	Delete(id string) error
	Update(id, title, desc string) (*Task, error)
	BulkImport([]*Task)
}
