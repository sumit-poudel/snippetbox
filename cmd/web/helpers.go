package main

import (
	"net/http"
	"runtime/debug"
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
