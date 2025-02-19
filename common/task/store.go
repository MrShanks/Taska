package task

type Store interface {
	GetTasks() []*Task
	New(*Task)
}
