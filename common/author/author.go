package author

import (
	"github.com/google/uuid"
)

type Author struct {
	ID        uuid.UUID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"Lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}
