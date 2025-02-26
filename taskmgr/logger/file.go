package logger

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type FileTransactionLogger struct {
	events chan<- Event
	errors <-chan error
	lastID int
	file   *os.File
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("couldn't open transaction log file: %v", err)
	}

	return &FileTransactionLogger{file: file}, nil
}

func (l *FileTransactionLogger) WriteNew(id uuid.UUID, title, desc string) {
	l.events <- Event{Type: New, TaskID: id, TaskTitle: title, TaskDesc: desc}
}

func (l *FileTransactionLogger) WriteMod(id uuid.UUID, title, desc string) {
	l.events <- Event{Type: Mod, TaskID: id, TaskTitle: title, TaskDesc: desc}
}

func (l *FileTransactionLogger) WriteDel(id uuid.UUID) {
	l.events <- Event{Type: Del, TaskID: id}
}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)

	outEvent := make(chan Event)
	outError := make(chan error)

	go func() {
		var err error
		var e Event
		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			if e.ID, err = strconv.Atoi(strings.Split(line, "\t")[0]); err != nil {
				outError <- fmt.Errorf("couldn't parse line %q: %v", line, err)
			}

			if l.lastID >= e.ID {
				outError <- fmt.Errorf("transaction number out of sequence")
				return
			}

			l.lastID = e.ID

			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("couldn't scan transaction log file: %v", err)
			return
		}
	}()

	return outEvent, outError
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) Run() {
	events := make(chan Event)
	l.events = events
	errors := make(chan error)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastID++

			_, err := fmt.Fprintf(
				l.file,
				"%d\t%s\t%v\t%-20s\t%s\n",
				l.lastID, e.Type, e.TaskID, e.TaskTitle, e.TaskDesc)

			if err != nil {
				errors <- err
				return
			}
		}
	}()
}
