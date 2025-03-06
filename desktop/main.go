package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func StartUp() {
	Taska := app.New()
	main := Taska.NewWindow("Taska")

	main.Resize(fyne.NewSize(1024, 768))

	main.SetContent(createGUI())
	main.SetMaster()
	main.ShowAndRun()
}

func main() {
	StartUp()
}
