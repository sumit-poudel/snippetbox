package main

import (
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/sumit-poudel/snippetbox/internal/models"
)

type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "Port address for the server")
	dsn := flag.String("dsn", "host=localhost user=web password=pass dbname=snippetbox sslmode=disable", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, errr := openDB(*dsn)

	if errr != nil {
		log.Fatal(errr)
	}
	defer db.Close()

	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: slog.NewLogLogger(logger.Handler(),slog.LevelError),
		Handler:  app.routes(),
	}
	logger.Info("Starting server ", slog.String("address", *addr))
	err := srv.ListenAndServe() // <-- this takes a handler
	logger.Error(err.Error())
	os.Exit(1)
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
