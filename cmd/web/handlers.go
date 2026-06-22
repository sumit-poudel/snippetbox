package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sumit-poudel/snippetbox/internal/models"
)

func (app *application) home(res http.ResponseWriter, req *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			res.Write([]byte("No snippets to get"))
			return
		} else {
			app.serverError(res, req, err)
			return
		}
	}
	data := app.newTemplateDate(req)
	data.Snippets = snippets
	app.render(res, req, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil || id < 0 {
		app.notFound(res, req)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			app.notFound(res, req)
			return
		} else {
			app.serverError(res, req, err)
		}
	}
	data := app.newTemplateDate(req)
	data.Snippet = snippet
	app.render(res, req, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("create a snippet"))
}

func (app *application) snippetCreatePost(res http.ResponseWriter, req *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(res, req, err)
		return
	}
	http.Redirect(res, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
