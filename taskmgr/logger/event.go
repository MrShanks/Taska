package logger

import "github.com/google/uuid"

type Operation string

const (
	New Operation = "NEW"
	Get Operation = "GET"
	Mod Operation = "MOD"
	Del Operation = "DEL"
)

type Event struct {
	ID        int
	Type      Operation
	TaskID    uuid.UUID
	TaskTitle string
	TaskDesc  string
}
