package storage

import (
	"log"

	"github.com/google/uuid"

	"github.com/MrShanks/Taska/common/task"
)

// InMemoryDatabase implements the taskStore interface
type InMemoryDatabase struct {
	Tasks map[uuid.UUID]*task.Task
}

func (imd *InMemoryDatabase) GetTasks() map[uuid.UUID]*task.Task {
	return imd.Tasks
}

func (imd *InMemoryDatabase) New(task *task.Task) uuid.UUID {
	id := uuid.New()
	imd.Tasks[id] = task
	return id
}

func (imd *InMemoryDatabase) Delete(id string) {
	UUID := uuid.MustParse(id)
	if _, ok := imd.Tasks[UUID]; ok {
		delete(imd.Tasks, UUID)
		log.Printf("Task with ID: %v has been deleted", UUID)
	} else {
		log.Printf("Couldn't find a task with that ID: %v", UUID)
	}
}
