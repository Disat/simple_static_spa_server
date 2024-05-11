package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type config struct {
	addr    string
	dirRoot string
}
type application struct {
	// errorLog *log.Logger
	infoLog *log.Logger
	config  config
}

func main() {
	var app application

	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	flag.StringVar(&app.config.addr, "port", ":8181", "static SPA(simgle page application) server listen port like :8080, default is :8181")
	flag.StringVar(&app.config.dirRoot, "root directory", "./dist", "root of static file directory in server like ./public, default is ./frontend")

	flag.Parse()

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(app.customHandler))

	server := &http.Server{
		Handler: app.logRequest(mux),
		Addr:    app.config.addr,
	}

	app.infoLog.Printf("The static SPA server listen on %s\n", app.config.addr)
	app.infoLog.Printf("the server static root path is %s\n", app.config.dirRoot)

	server.ListenAndServe()
}

func (app *application) customHandler(w http.ResponseWriter, r *http.Request) {

	_, err := os.Stat(filepath.Join(app.config.dirRoot, r.URL.Path))
	if err != nil {

		http.ServeFile(w, r, app.config.dirRoot+"/index.html")
		return
	}

	http.FileServer(http.Dir(app.config.dirRoot)).ServeHTTP(w, r)
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}
