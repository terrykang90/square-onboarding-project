package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	templates *template.Template
}

func main() {
	addr := flag.String("addr", ":4100", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		errorLog.Fatal("Error loading .env file")
	}

	var templateFunc = template.FuncMap{
		"add":           add,
		"remove_quotes": remove_quotes,
	}

	app := &application{
		errorLog:  errorLog,
		infoLog:   infoLog,
		templates: template.Must(template.New("").Funcs(templateFunc).ParseGlob("ui/html/*")),
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
