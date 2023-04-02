package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func newProjectHandler(w http.ResponseWriter, r *http.Request) {
	cl := getFlaarumClient()

	err := cl.CreateProject(r.FormValue("name"))
	if err != nil {
		ErrorPage(w, err)
		return
	}

	http.Redirect(w, r, "/?project="+r.FormValue("name"), http.StatusTemporaryRedirect)

}

func newTableHandler(w http.ResponseWriter, r *http.Request) {
	cl := getFlaarumClient()
	cl.ProjName = r.FormValue("current_project")

	err := cl.CreateTable(r.FormValue("stmt"))
	if err != nil {
		ErrorPage(w, err)
		return
	}

	http.Redirect(w, r, "/?project="+r.FormValue("current_project"), http.StatusTemporaryRedirect)
}

func loadTableHandler(w http.ResponseWriter, r *http.Request) {
	cl := getFlaarumClient()
	cl.ProjName = r.FormValue("project")

	tableName := r.FormValue("table")

	rows, err := cl.Search(fmt.Sprintf(`
		table: %s
		limit: 100
		order_by: id asc
	`, tableName))
	if err != nil {
		ErrorPage(w, err)
		return
	}

	vnum, _ := cl.GetCurrentTableVersionNum(tableName)
	tableDefnParsed, _ := cl.GetTableStructureParsed(tableName, vnum)

	fields := make([]string, 0)
	innerFields := make([]string, 0)
	for _, fieldStruct := range tableDefnParsed.Fields {
		fields = append(fields, fieldStruct.FieldName+"["+fieldStruct.FieldType+"]")
		innerFields = append(innerFields, fieldStruct.FieldName)
	}

	retRows := make([][]any, 0)
	for _, row := range *rows {
		reportedRowSlice := make([]any, 0)
		for _, field := range innerFields {
			reportedRowSlice = append(reportedRowSlice, row[field])
		}
		retRows = append(retRows, reportedRowSlice)
	}

	count, _ := cl.AllRowsCount(tableName)
	type Context struct {
		Table          string
		Fields         []string
		Rows           [][]any
		AllRowsCount   int64
		CurrentVersion int
	}
	tmpl := template.Must(template.ParseFS(content, "templates/table_view.html"))
	tmpl.Execute(w, Context{tableName, fields, retRows, count, int(vnum)})
}
