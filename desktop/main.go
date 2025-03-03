package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/MrShanks/Taska/common/task"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Form Layout")

	titleLabel := widget.NewLabel("Title:")
	titleInput := widget.NewEntry()
	descLabel := widget.NewLabel("Desc:")
	descInput := widget.NewEntry()

	submit := widget.NewButton("Submit", func() {

		t := task.New(titleInput.Text, descInput.Text)

		jsonTask, err := json.Marshal(t)
		if err != nil {
			fmt.Println("Error:", err)
		}

		http.Post("http://localhost:8080/new", "application/json", bytes.NewBuffer(jsonTask))
	})

	grid := container.New(layout.NewFormLayout(), titleLabel, titleInput, descLabel, descInput)

	form := container.NewVBox(grid, submit)

	myWindow.SetContent(form)
	myWindow.ShowAndRun()
}

