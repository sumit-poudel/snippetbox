package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sumit-poudel/snippetbox/internal/models"
	"github.com/sumit-poudel/snippetbox/internal/validator"
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
	data := app.newTemplateData(req)
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
	data := app.newTemplateData(req)
	data.Snippet = snippet
	app.render(res, req, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(res http.ResponseWriter, req *http.Request) {
	date := app.newTemplateData(req)
	date.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(res, req, http.StatusOK, "create.tmpl.html", date)
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) snippetCreatePost(res http.ResponseWriter, req *http.Request) {

	var form snippetCreateForm
	err := app.decodePostForm(req, &form)
	if err != nil {
		app.clientError(res, http.StatusBadRequest)
		return
	}

	err = app.formDecoder.Decode(&form, req.PostForm)

	if err != nil {
		app.clientError(res, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more then 100 character long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1,7,365")

	if !form.Valid() {
		data := app.newTemplateData(req)
		data.Form = form
		app.render(res, req, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(res, req, err)
		return
	}

	app.sessionManager.Put(req.Context(), "flash", "Snippet successfully created!")
	http.Redirect(res, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
