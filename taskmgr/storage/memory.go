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

func (imd *InMemoryDatabase) GetTasks() map[uuid.UUID]*task.Task {
	return imd.Tasks
}

func (imd *InMemoryDatabase) New(task *task.Task) uuid.UUID {
	task.ID = uuid.New()
	imd.Tasks[task.ID] = task
	return task.ID
}

func (imd *InMemoryDatabase) Update(id, title, desc string) (*task.Task, error) {
	updateTask := imd.Tasks[uuid.MustParse(id)]
	if _, ok := imd.Tasks[updateTask.ID]; !ok {
		return nil, fmt.Errorf("task not found")
	}

	if title != "" {
		updateTask.Title = title
	}

	if desc != "" {
		updateTask.Desc = desc
	}

	return updateTask, nil
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
