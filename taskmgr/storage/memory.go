package storage

import (
	"github.com/MrShanks/Taska/common/task"
)

// inMemoryDatabase implements the taskStore interface
type InMemoryDatabase struct {
	Tasks []*task.Task
}

func (imd *InMemoryDatabase) GetTasks() []*task.Task {
	return imd.Tasks
}

func (imd *InMemoryDatabase) New(task *task.Task) {
	imd.Tasks = append(imd.Tasks, task)
}
