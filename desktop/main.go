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

func main() {
	Taska := app.New()
	mainWindow := Taska.NewWindow("Form Layout")

	titleLabel := widget.NewLabel("Title:")
	titleInput := widget.NewEntry()
	descLabel := widget.NewLabel("Desc:")
	descInput := widget.NewEntry()

	tasks := make(map[uuid.UUID]*task.Task)

	submit := widget.NewButton("Submit", func() {

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
	})

	tasksContainer := container.NewVBox()

	get := widget.NewButton("Get", func() {
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

		if err := json.Unmarshal(data, &tasks); err != nil {
			fmt.Println("Error: ", err)
		}

		tasksContainer.Objects = nil

		for id, task := range tasks {
			tasksContainer.Add(container.NewHBox(widget.NewLabel(id.String()+": "), widget.NewLabel(task.Title)))

		}

		tasksContainer.Refresh()
	})

	grid := container.New(layout.NewFormLayout(), titleLabel, titleInput, descLabel, descInput)

	form := container.NewVBox(grid, submit, get, tasksContainer)

	mainWindow.Resize(fyne.NewSize(600, 500))
	mainWindow.SetContent(form)

	mainWindow.ShowAndRun()
}
