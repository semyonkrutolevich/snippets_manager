package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"krutolevichsemyon.life/snippetbox/pkg/models/psql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *psql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// TODO Change to enviroment variables
	DB_HOST := "127.0.0.1"
	DB_PORT := "5432"
	DB_USER := "snippets_user"
	DB_PASSWORD := "snippets_pass"
	DB_NAME := "snippets_db"

	psqlDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, err := openDB(psqlDSN)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("../client/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		snippets:      &psql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	host := flag.String("h", "127.0.0.1:5000", "Network Interface HTTP")
	flag.Parse()

	server := &http.Server{
		Addr:     *host,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting web server on %s", *host)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
