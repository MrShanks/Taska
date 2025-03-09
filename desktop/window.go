package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/google/uuid"
)

func AddTaskToUI(tasksContainer *fyne.Container, title, desc string, id uuid.UUID) {
	taskCard := widget.NewCard(title, "", nil)
	descLabel := widget.NewLabel(desc)

	deleteIcon := widget.NewButtonWithIcon("delete", theme.DeleteIcon(), DeleteTask(id, tasksContainer))
	editIcon := widget.NewButtonWithIcon("edit", theme.DocumentIcon(), UpdateTask(id, tasksContainer))

	descRow := container.NewHBox(descLabel, layout.NewSpacer(), editIcon, deleteIcon)

	taskBox := container.NewVBox(
		taskCard,
		descRow,
		widget.NewSeparator(),
	)

	tasksContainer.Add(taskBox)
}

func createGUI() *fyne.Container {

	titleLabel := widget.NewLabel("Title:")
	titleInput := widget.NewEntry()

	descLabel := widget.NewLabel("Desc:")
	descInput := widget.NewMultiLineEntry()

	form := container.New(
		layout.NewFormLayout(),
		titleLabel,
		titleInput,
		descLabel,
		descInput)

	tasksContainer := container.NewVBox()
	GetTasks(tasksContainer)

	submitBtn := widget.NewButton(
		"Submit",
		SubmitNewTask(titleInput, descInput, tasksContainer))

	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.SearchIcon(), func() {
		// TODO: implement search
	}))

	left := container.NewVBox(
		container.NewPadded(form),
		container.NewPadded(submitBtn),
		layout.NewSpacer(),
	)

	dividers := [2]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(),
	}

	objs := []fyne.CanvasObject{toolbar, left, tasksContainer, dividers[0], dividers[1]}
	mainLayout := container.New(newTaskaLayout(toolbar, left, tasksContainer, dividers), objs...)

	return mainLayout
}
