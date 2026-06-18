package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(res http.ResponseWriter, req *http.Request) { // <-- handler to handel response
	if req.URL.Path != "/" {
		app.notFound(res, *req)
		return
	}
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partial/navbar.tmpl.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(res, err)
		return
	}
	error := ts.ExecuteTemplate(res, "base", nil)
	if error != nil {
		app.serverError(res, err)
	}
}

func (app *application) snippetView(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		res.Header().Set("Allow", http.MethodGet)
		app.clientError(res, 405)
		return
	}
	qry, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || qry < 0 {
		app.notFound(res, *req)
		return
	}
	res.Header().Set("content-type", "application/json")         // <-- this does case sanitation
	res.Header()["X-XSS-PROTECTION"] = []string{"1; mode=block"} // <- if you dont want case sanitation
	fmt.Fprintf(res, "The number you query is %d...", qry)
}

func (app *application) snippetCreate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.Header().Set("Allow", http.MethodPost) // <-- must be called before .Write for custom status code
		res.WriteHeader(http.StatusMethodNotAllowed)
		res.Write([]byte("Method not allowed"))
		return
	}
	res.Write([]byte("Snippet created...")) // <- only this will send 200 status code
}
