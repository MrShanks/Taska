package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/MrShanks/Taska/common/task"
	"github.com/google/uuid"
)

func SubmitNewTask(titleInput, descInput *widget.Entry) func() {
	return func() {
		t := task.New(titleInput.Text, descInput.Text)

		jsonTask, err := json.Marshal(t)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		data := bytes.NewBuffer(jsonTask)

		response, err := http.Post("http://localhost:8080/new", "application/json", data)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer response.Body.Close()
	}
}

func GetTasks(ctr *fyne.Container) func() {
	return func() {

		response, err := http.Get("http://localhost:8080/tasks")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		tasks := make(map[uuid.UUID]*task.Task)
		if err := json.Unmarshal(data, &tasks); err != nil {
			fmt.Println("Error: ", err)
		}

		ctr.Objects = nil

		for id, task := range tasks {
			ctr.Add(container.NewHBox(widget.NewLabel(id.String()+": "), widget.NewLabel(task.Title)))

		}

		ctr.Refresh()
	}
}

func SetupWindowContent() *fyne.Container {

	titleLabel := widget.NewLabel("Title:")
	titleInput := widget.NewEntry()

	descLabel := widget.NewLabel("Desc:")
	descInput := widget.NewEntry()

	submit := widget.NewButton("Submit", SubmitNewTask(titleInput, descInput))

	tasksContainer := container.NewVBox()
	get := widget.NewButton("Get", GetTasks(tasksContainer))

	form := container.New(layout.NewFormLayout(), titleLabel, titleInput, descLabel, descInput)

	return container.NewVBox(form, submit, get, tasksContainer)
}

func main() {
	Taska := app.New()
	mainWindow := Taska.NewWindow("Form Layout")

	content := SetupWindowContent()

	mainWindow.Resize(fyne.NewSize(600, 500))
	mainWindow.SetContent(content)

	mainWindow.ShowAndRun()
}
