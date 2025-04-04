package storage

import (
	"fmt"
	"log"
	"sort"

	"github.com/google/uuid"

	"github.com/MrShanks/Taska/common/task"
)

// InMemoryDatabase implements the taskStore interface
type InMemoryDatabase struct {
	Tasks map[uuid.UUID]*task.Task
}

func (imd *InMemoryDatabase) Search(keyword, authorID string) ([]*task.Task, error) {
	return nil, nil
}

func (imd *InMemoryDatabase) GetOne(id, authorID string) (*task.Task, error) {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %s", id)
	}

	t, ok := imd.Tasks[UUID]
	if !ok {
		return nil, fmt.Errorf("task with ID: %v not found", UUID)
	}

	return t, nil
}

func (imd *InMemoryDatabase) GetTasks(authorID string) []*task.Task {
	tasks := []*task.Task{}

	for _, val := range imd.Tasks {
		tasks = append(tasks, val)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Title < tasks[j].Title
	})

	return tasks
}

func (imd *InMemoryDatabase) New(task *task.Task) uuid.UUID {
	task.ID = uuid.New()
	imd.Tasks[task.ID] = task
	return task.ID
}

func (imd *InMemoryDatabase) Update(id, title, desc, author string) (*task.Task, error) {
	taskID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid")
	}

	if _, ok := imd.Tasks[taskID]; !ok {
		return nil, fmt.Errorf("task not found")
	}

	updateTask := imd.Tasks[uuid.MustParse(id)]

	if title != "" {
		updateTask.Title = title
	}

	if desc != "" {
		updateTask.Desc = desc
	}

	return updateTask, nil
}

func (imd *InMemoryDatabase) Delete(id, author string) error {
	UUID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid uuid: %s", id)
	}

	if _, ok := imd.Tasks[UUID]; !ok {
		return fmt.Errorf("task with ID: %v not found", UUID)
	}

	delete(imd.Tasks, UUID)
	log.Printf("Task with ID: %v has been deleted", UUID)
	return nil
}

func (imd *InMemoryDatabase) BulkImport(tasks []*task.Task, author string) {
	for _, task := range tasks {
		imd.Tasks[task.ID] = task
	}
}
