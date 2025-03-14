package author

import (
	"github.com/MrShanks/Taska/common/task"
	"github.com/google/uuid"
)

type Author struct {
	ID        uuid.UUID
	Firstname string
	Lastname  string
	Email     string
	Password  string
	Tasks     []*task.Task
}
