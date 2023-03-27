package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

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
	os.Setenv("FYNE_SCALE", "0.9")

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

	doRightRefreshChan := make(chan string)
	doLeftRefreshChan := make(chan bool)

	mainContent := container.NewVBox()

	getTableUI := func(projectName, tableName string) *fyne.Container {
		vnum, _ := cl.GetCurrentTableVersionNum(tableName)
		versionLabel := widget.NewLabel(fmt.Sprintf("Version: %d", vnum))
		topTableBar := container.NewHBox()
		utsBtn := widget.NewButton("Update Table Structure", func() {
			l := widget.NewLabel("statement")
			statementEntry := widget.NewMultiLineEntry()
			vnum, _ := cl.GetCurrentTableVersionNum(tableName)
			tableStructStmt, _ := cl.GetTableStructure(tableName, vnum)
			statementEntry.SetText(tableStructStmt)

			callbackFunc := func(b bool) {
				if b {
					cl.ProjName = projectName
					err = cl.UpdateTableStructure(statementEntry.Text)
					if err != nil {
						dialog.ShowInformation("Error creating table", err.Error(), myWindow)
					}
					doLeftRefreshChan <- true
				}
			}
			dialogContent := container.New(&diags{}, l, statementEntry)
			dialogContent.Resize(fyne.NewSize(600, 400))

			dialog.ShowCustomConfirm("Update Table Structure", "Update", "Cancel", dialogContent, callbackFunc, myWindow)
		})

		topTableBar.Add(versionLabel)
		topTableBar.Add(utsBtn)

		deleteTableBtn := widget.NewButton("Delete Table", func() {
			deleteTableCallback := func(b bool) {
				if b {
					tables, _ := cl.ListTables()

					err = cl.DeleteTable(tableName)
					if err != nil {
						dialog.ShowInformation("Error creating table", err.Error(), myWindow)
					}
					doLeftRefreshChan <- true
					if len(tables) > 0 {
						doRightRefreshChan <- tables[0]
					} else {
						doRightRefreshChan <- ""
					}
				}
			}
			dialog.ShowConfirm("Delete Table "+tableName+" Confirmation", "Do you really want to delete this table",
				deleteTableCallback, myWindow)
		})
		topTableBar.Add(deleteTableBtn)

		tableContent := container.NewVBox(topTableBar, widget.NewSeparator(), widget.NewLabel(tableName))
		return tableContent
	}

	loadLeftUI := func(project string) *fyne.Container {
		titleLabel := widget.NewLabel("Tables")
		cl.ProjName = project
		tables, _ := cl.ListTables()

		createTableBtn := widget.NewButton("Create Table", func() {
			l := widget.NewLabel("statement")
			statementEntry := widget.NewMultiLineEntry()
			trueCreateTableFunc := func(b bool) {
				if b {
					cl.ProjName = project
					err = cl.CreateTable(statementEntry.Text)
					if err != nil {
						dialog.ShowInformation("Error creating table", err.Error(), myWindow)
					}
					doLeftRefreshChan <- true
				}
			}
			dialogContent := container.New(&diags{}, l, statementEntry)
			dialogContent.Resize(fyne.NewSize(600, 400))

			dialog.ShowCustomConfirm("New Table", "Create", "Cancel", dialogContent, trueCreateTableFunc, myWindow)
		})

		UIContent := container.NewVBox(
			container.NewHBox(titleLabel, createTableBtn),
		)

		for _, tableName := range tables {
			makeTableFunc := func(tableName string) func() {
				return func() {
					doRightRefreshChan <- tableName
				}
			}
			tableBtn := widget.NewButton(tableName, makeTableFunc(tableName))
			UIContent.Add(tableBtn)
		}

		return UIContent
	}

	leftContent := container.NewVBox()
	projectsSwitch := widget.NewSelect(projects, func(s string) {
		content := loadLeftUI(s)
		leftContent.RemoveAll()
		leftContent.Add(content)
		leftContent.Refresh()
	})

	projectsSwitch.SetSelected("first_proj")
	// refresh left UI thread
	go func() {
		for {
			<-doLeftRefreshChan

			content := loadLeftUI(projectsSwitch.Selected)
			leftContent.RemoveAll()
			leftContent.Add(content)
			leftContent.Refresh()

			time.Sleep(time.Second)
		}
	}()

	// refresh left UI thread
	go func() {
		for {
			tableName := <-doRightRefreshChan

			if tableName == "" {
				mainContent.RemoveAll()
			} else {
				content := getTableUI(projectsSwitch.Selected, tableName)
				mainContent.RemoveAll()
				mainContent.Add(content)
				mainContent.Refresh()
			}

			time.Sleep(time.Second)
		}
	}()

	newProjectBtn := widget.NewButton("New Project", func() {
		content := make([]*widget.FormItem, 0)
		content = append(content, widget.NewFormItem("name", widget.NewEntry()))
		dialog.ShowForm("New Project", "Create", "Cancel", content, func(b bool) {
			if b {
				inputs := getFormInputs(content)
				cl.CreateProject(inputs["name"])
				content := loadLeftUI(inputs["name"])
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

	windowContent := container.NewBorder(topBar, nil, leftContent, nil, mainContent)

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
