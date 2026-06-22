package main

import (
	"fmt"
	"net/http"
)

func commonHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		res.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		res.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		res.Header().Set("X-Content-Type-Options", "nosniff")
		res.Header().Set("X-Frame-Options", "deny")
		res.Header().Set("X-XSS-Protection", "0")
		res.Header().Set("Server", "Go")
		next.ServeHTTP(res, req)
	})
}

func (app *application) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var (
			ip     = req.RemoteAddr
			proto  = req.Proto
			method = req.Method
			uri    = req.URL.RequestURI()
		)
		app.logger.Info("Received requeest", "ip", ip, "protocol", proto, "method", method, "uri", uri)
		next.ServeHTTP(res, req)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				res.Header().Set("Connection", "close")
				app.serverError(res, req, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(res, req)
	})
}
