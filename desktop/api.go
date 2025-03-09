package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/MrShanks/Taska/common/task"
	"github.com/google/uuid"
)

func SubmitNewTask(titleInput, descInput *widget.Entry, ctr *fyne.Container) func() {
	return func() {
		t := task.New(titleInput.Text, descInput.Text)

		jsonTask, err := json.Marshal(t)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		data := bytes.NewBuffer(jsonTask)

		request, err := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8080/new", data)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		client := http.Client{}

		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		response.Body.Close()

		GetTasks(ctr)
	}
}

func GetTasks(ctr *fyne.Container) {
	request, err := http.NewRequestWithContext(context.Background(), "GET", "http://localhost:8080/tasks", nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	client := http.Client{}
	response, err := client.Do(request)
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

	for _, task := range tasks {
		AddTaskToUI(ctr, task.Title, task.Desc, task.ID)
	}

	ctr.Refresh()
}

func DeleteTask(id uuid.UUID, ctr *fyne.Container) func() {
	return func() {
		request, err := http.NewRequestWithContext(context.Background(), "DELETE", fmt.Sprintf("http://localhost:8080/delete/%s", id.String()), nil)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		client := http.Client{}

		response, err := client.Do(request)
		if err != nil {
			fmt.Printf("Couldn't get a response from the server: %v\n", err)
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotFound {
			fmt.Println("Task not found")
		}

		GetTasks(ctr)
	}
}

func UpdateTask(id uuid.UUID, ctr *fyne.Container, oldTitle, oldDesc string) func() {
	return func() {
		editWindow := fyne.CurrentApp().NewWindow("Edit")

		newTitle := widget.NewEntry()
		newTitle.SetPlaceHolder(oldTitle)
		newDesc := widget.NewMultiLineEntry()
		newDesc.SetPlaceHolder(oldDesc)

		submitBtn := widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
			body := &bytes.Buffer{}
			body.Write(fmt.Appendf(nil, `{"title":"%s","desc":"%s"}`, newTitle.Text, newDesc.Text))

			request, err := http.NewRequestWithContext(context.Background(), "PUT", fmt.Sprintf("http://localhost:8080/update/%s", id.String()), body)
			if err != nil {
				fmt.Println("Error: ", err)
			}

			client := http.Client{}

			response, err := client.Do(request)
			if err != nil {
				fmt.Printf("Couldn't get a response from the server: %v\n", err)
			}
			defer response.Body.Close()

			editWindow.Close()
			GetTasks(ctr)
		})

		editWindow.SetContent(container.NewVBox(newTitle, newDesc, submitBtn))

		editWindow.Resize(fyne.NewSize(400, 200))

		editWindow.Show()
	}
}
