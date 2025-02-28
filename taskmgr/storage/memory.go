package storage

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/MrShanks/Taska/common/task"
)

// InMemoryDatabase implements the taskStore interface
type InMemoryDatabase struct {
	Tasks map[uuid.UUID]*task.Task
}

func (imd *InMemoryDatabase) GetOne(id string) (*task.Task, error) {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s", id)
	}

	if _, ok := imd.Tasks[UUID]; !ok {
		return nil, fmt.Errorf("task with ID: %v not found", UUID)
	}

	return imd.Tasks[UUID], nil
}

func (imd *InMemoryDatabase) GetTasks() map[uuid.UUID]*task.Task {
	return imd.Tasks
}

func (imd *InMemoryDatabase) New(task *task.Task) uuid.UUID {
	task.ID = uuid.New()
	imd.Tasks[task.ID] = task
	return task.ID
}

func (imd *InMemoryDatabase) Delete(id string) error {
	UUID := uuid.MustParse(id)
	if _, ok := imd.Tasks[UUID]; !ok {
		return fmt.Errorf("task with ID: %v not found", UUID)
	}

	delete(imd.Tasks, UUID)
	log.Printf("Task with ID: %v has been deleted", UUID)
	return nil
}
