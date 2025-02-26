package logger

import "github.com/google/uuid"

type TransactionLogger interface {
	WriteNew(id uuid.UUID, title, desc string)
	WriteMod(id uuid.UUID, title, desc string)
	WriteDel(id uuid.UUID)

	ReadEvents()(<-chan Event,  <-chan error)

	Err() <-chan error
	Run()
}
