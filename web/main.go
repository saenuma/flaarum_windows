package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/gs/{obj}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rawObj, err := contentStatics.ReadFile("statics/" + vars["obj"])
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Disposition", "attachment; filename="+vars["obj"])
		contentType := http.DetectContentType(rawObj)
		w.Header().Set("Content-Type", contentType)
		w.Write(rawObj)
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFS(content, "templates/app.html"))
		tmpl.Execute(w, nil)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
