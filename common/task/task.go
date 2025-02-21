package task

type Task struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func New(title, desc string) *Task {
	return &Task{
		Title: title,
		Desc:  desc,
	}
}
