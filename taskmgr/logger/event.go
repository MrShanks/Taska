package logger

import (
	"time"

	"github.com/google/uuid"
)

type Operation string

const (
	New Operation = "NEW"
	Mod Operation = "MOD"
	Del Operation = "DEL"
)

type Event struct {
	ID        int
	Type      Operation
	TaskID    uuid.UUID
	TaskTitle string
	TaskDesc  string
	Time      time.Time
}
