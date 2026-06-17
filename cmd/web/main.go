package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux() //<-- multiplexr to handel the respose
	fileserver := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/",http.StripPrefix("/static",fileserver))
	mux.HandleFunc("/", home) // <-- router
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	log.Println("Starting server at 4000")
	err := http.ListenAndServe(":4000", mux) // <-- this takes a handler
	log.Fatal(err)
}
/*
listen and serve takes a handler (http handler)
cons:
	have to use conditons for responses
	too much complicated for other then simple server

So use a multiplexer that maps the route with handler or handler function
handler and handler fuc are diff
handler inherits the serveHTTP interface
where as handler func is just a func which can be changed to handler using handlerFunc()
*/
