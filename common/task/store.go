package task

import "github.com/google/uuid"

type Store interface {
	GetTasks(authorID string) []*Task
	GetOne(id, authorID string) (*Task, error)
	New(task *Task) uuid.UUID
	Delete(id, authorID string) error
	Update(id, title, desc, authorID string) (*Task, error)
	BulkImport(tasks []*Task, authorID string)
	Search(keyword, authorID string) ([]*Task, error)
}
