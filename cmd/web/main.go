package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/sumit-poudel/snippetbox/internal/models"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *models.SnippetModel
}

func main() {
	// declare flag for the runtime
	addr := flag.String("addr", ":4000", "Port address for the server")
	dsn := flag.String("dsn", "host=localhost user=web password=pass dbname=snippetbox sslmode=disable", "MySQL data source name")
	// parse the flag
	flag.Parse()

	// custom log messages and levels
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile)

	db, errr := openDB(*dsn)

	if errr != nil {
		log.Fatal(errr)
	}
	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errLog,
		snippets: &models.SnippetModel{DB: db},
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

func openDB(dns string) (*sql.DB, error) {
	db, errr := sql.Open("postgres", dns)
	if errr != nil {
		return nil, errr
	}
	if errr = db.Ping(); errr != nil {
		return nil, errr
	}

	return db, errr
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
