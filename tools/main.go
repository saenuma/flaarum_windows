package main

import (
	"fmt"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/saenuma/flaarum"
	"github.com/saenuma/flaarum/flaarum_shared"
)

func main() {
	os.Setenv("FYNE_THEME", "dark")

	myApp := app.New()
	myWindow := myApp.NewWindow("hanan: a more comfortable shell / terminal")

	var keyStr string
	inProd := flaarum_shared.GetSetting("in_production")
	if inProd == "" {
		fmt.Println("unexpected error. Have you installed  and launched flaarum?")
		os.Exit(1)
	}
	if inProd == "true" {
		keyStrPath := flaarum_shared.GetKeyStrPath()
		raw, err := os.ReadFile(keyStrPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		keyStr = string(raw)
	} else {
		keyStr = "not-yet-set"
	}
	port := flaarum_shared.GetSetting("port")
	if port == "" {
		fmt.Println("unexpected error. Have you installed  and launched flaarum?")
		os.Exit(1)
	}
	var cl flaarum.Client

	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("Invalid port setting.")
		os.Exit(1)
	}

	if portInt != flaarum_shared.PORT {
		cl = flaarum.NewClientCustomPort("127.0.0.1", keyStr, "first_proj", portInt)
	} else {
		cl = flaarum.NewClient("127.0.0.1", keyStr, "first_proj")
	}

	err = cl.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	projects, err := cl.ListProjects()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loadUI := func(project string) *fyne.Container {
		titleLabel := widget.NewLabel("Tables")
		cl.ProjName = project
		tables, _ := cl.ListTables()

		createTableBtn := widget.NewButton("Create Table", func() {

		})

		UIContent := container.NewVBox(
			container.NewHBox(titleLabel, createTableBtn),
		)
		for _, tableName := range tables {
			UIContent.Add(widget.NewButton(tableName, func() {

			}))
		}

		return UIContent
	}

	leftContent := container.NewVBox()
	projectsSwitch := widget.NewSelect(projects, func(s string) {
		content := loadUI(s)
		leftContent.RemoveAll()
		leftContent.Add(content)
		leftContent.Refresh()
	})

	newProjectBtn := widget.NewButton("New Project", func() {
		content := make([]*widget.FormItem, 0)
		content = append(content, widget.NewFormItem("name", widget.NewEntry()))
		dialog.ShowForm("New Project", "Create", "Cancel", content, func(b bool) {
			if b {
				inputs := getFormInputs(content)
				cl.CreateProject(inputs["name"])
				content := loadUI(inputs["name"])
				leftContent.RemoveAll()
				leftContent.Add(content)
				leftContent.Refresh()
			}
		}, myWindow)
	})

	topBar := container.NewVBox(
		container.NewCenter(container.NewHBox(projectsSwitch, newProjectBtn)),
		widget.NewSeparator(),
	)

	windowContent := container.NewBorder(topBar, nil, leftContent, nil, nil)

	myWindow.SetContent(windowContent)
	myWindow.Resize(fyne.NewSize(1200, 600))

	myWindow.ShowAndRun()
}

func getFormInputs(content []*widget.FormItem) map[string]string {
	inputs := make(map[string]string)
	for _, formItem := range content {
		entryWidget, ok := formItem.Widget.(*widget.Entry)
		if ok {
			inputs[formItem.Text] = entryWidget.Text
			continue
		}

		selectWidget, ok := formItem.Widget.(*widget.Select)
		if ok {
			inputs[formItem.Text] = selectWidget.Selected
		}
	}

	return inputs
}
