package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	// declare flag for the runtime
	addr := flag.String("addr", ":4000", "Port address for the server")
	// parse the flag
	flag.Parse()

	// custom log messages and levels
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server at %s", *addr)
	err := srv.ListenAndServe() // <-- this takes a handler
	errLog.Fatal(err)
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
