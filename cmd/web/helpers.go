package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form"
)

func (app *application) serverError(res http.ResponseWriter, req *http.Request, err error) {
	var (
		method = req.Method
		uri    = req.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "URL", uri, "method", method, "trace", trace)
	http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(res http.ResponseWriter, status int) {
	http.Error(res, http.StatusText(status), status)
}

func (app *application) notFound(res http.ResponseWriter, req *http.Request) {
	http.NotFound(res, req)
}

func (app *application) render(res http.ResponseWriter, req *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("Template doesnt exist for %s", page)
		app.serverError(res, req, err)
	}
	buff := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buff, "base", data)
	if err != nil {
		app.serverError(res, req, err)
		return
	}
	res.WriteHeader(status)
	buff.WriteTo(res)
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear: time.Now().Year(),
		Flash: app.sessionManager.PopString(r.Context(),"flash"),
	}
}

func (app *application) decodePostForm(req *http.Request, dst any) error {
	err := req.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, req.PostForm)
	if err != nil {
		if _, ok := errors.AsType[*form.InvalidDecoderError](err); ok {
			panic(err)
		}
		return err
	}
	return err
}
