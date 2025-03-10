package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/google/uuid"
)

func populateToolbar() *widget.Toolbar {
	return widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), placeHolderFunction),
		widget.NewToolbarAction(theme.StorageIcon(), placeHolderFunction),
		widget.NewToolbarAction(theme.FileIcon(), placeHolderFunction),
		widget.NewToolbarAction(theme.DocumentPrintIcon(), placeHolderFunction),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), placeHolderFunction),
		widget.NewToolbarAction(theme.MailComposeIcon(), placeHolderFunction),
		widget.NewToolbarAction(theme.UploadIcon(), placeHolderFunction),
		widget.NewToolbarAction(theme.DownloadIcon(), placeHolderFunction),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), placeHolderFunction),
		widget.NewToolbarAction(theme.AccountIcon(), placeHolderFunction),
		widget.NewToolbarAction(theme.LoginIcon(), placeHolderFunction),
	)
}

func populateLeftContainer(tasks *fyne.Container) *fyne.Container {

	form := createTaskForm(tasks)

	leftContainer := container.NewVBox(
		container.NewPadded(widget.NewLabelWithStyle("New Task", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
		container.NewPadded(form),
		layout.NewSpacer(),
	)

	return leftContainer
}

func populateMainContent() *fyne.Container {
	mainContent := container.NewVBox()
	GetTasks(mainContent)

	return mainContent
}

func createTaskForm(tasks *fyne.Container) *fyne.Container {
	titleLabel := widget.NewLabel("Title:")
	titleInput := widget.NewEntry()

	descLabel := widget.NewLabel("Desc:")
	descInput := widget.NewMultiLineEntry()

	inputFields := container.New(
		layout.NewFormLayout(),
		titleLabel,
		titleInput,
		descLabel,
		descInput)

	submitBtn := widget.NewButton(
		"Submit",
		SubmitNewTask(titleInput, descInput, tasks))

	form := container.NewVBox(
		inputFields,
		submitBtn,
	)

	return form
}

func createGUI() *fyne.Container {

	toolbar := populateToolbar()

	mainContent := populateMainContent()

	leftContent := populateLeftContainer(mainContent)

	dividers := [2]fyne.CanvasObject{
		widget.NewSeparator(), widget.NewSeparator(),
	}

	objs := []fyne.CanvasObject{toolbar, leftContent, mainContent, dividers[0], dividers[1]}
	mainLayout := container.New(newTaskaLayout(toolbar, leftContent, mainContent, dividers), objs...)

	return mainLayout
}

func AddTaskToUI(tasksContainer *fyne.Container, title, desc string, id uuid.UUID) {
	taskCard := widget.NewCard(title, "", nil)
	descLabel := widget.NewLabel(desc)

	deleteIcon := widget.NewButtonWithIcon("delete", theme.DeleteIcon(), DeleteTask(id, tasksContainer))
	editIcon := widget.NewButtonWithIcon("edit", theme.DocumentIcon(), UpdateTask(id, tasksContainer, title, desc))

	descRow := container.NewHBox(descLabel, layout.NewSpacer(), editIcon, deleteIcon)

	taskBox := container.NewVBox(
		taskCard,
		descRow,
		widget.NewSeparator(),
	)

	tasksContainer.Add(taskBox)
}

func placeHolderFunction() {
	// this function is just a placeHolder for future button actions
}
