package logger

import "github.com/google/uuid"

type TransactionLogger interface {
	WriteNew(id uuid.UUID, title, desc string)
	WriteMod(id uuid.UUID, title, desc string)
	WriteDel(id uuid.UUID)

	Err() <-chan error
	Run()
}
