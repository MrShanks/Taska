package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func AddTaskToUI(tasksContainer *fyne.Container, title, desc string) {
	taskCard := widget.NewCard(title, desc, nil)
	taskCard.SetSubTitle(desc)

	taskBox := container.NewVBox(
		taskCard,
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

	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.SearchIcon(), func() {}))

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
