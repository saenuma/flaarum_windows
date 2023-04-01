package main

import "net/http"

func newProjectHandler(w http.ResponseWriter, r *http.Request) {
	cl := getFlaarumClient()

	err := cl.CreateProject(r.FormValue("name"))
	if err != nil {
		ErrorPage(w, err)
		return
	}

	http.Redirect(w, r, "/?project="+r.FormValue("name"), http.StatusTemporaryRedirect)

}
