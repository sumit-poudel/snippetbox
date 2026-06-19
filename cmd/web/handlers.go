package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(res http.ResponseWriter, req *http.Request) {

	res.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partial/navbar.tmpl.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, req, err)
		return
	}
	error := ts.ExecuteTemplate(res, "base", nil)
	if error != nil {
		app.serverError(res,req ,err)
	}
}

func (app *application) snippetView(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil || id < 0 {
		app.notFound(res, *req)
		return
	}
	res.Header().Set("content-type", "application/json")
	res.Header()["X-XSS-PROTECTION"] = []string{"1; mode=block"}
	fmt.Fprintf(res, "The number is %d...", id)
}

func (app *application) snippetCreate(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("create a snippet"))
}

func (app *application) snippetCreatePost(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Snippet created..."))
}
