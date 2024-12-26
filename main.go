package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed all:ui/tmpl
var TmplFiles embed.FS

//go:embed all:ui/static
var StaticFiles embed.FS

func main() {

	app := &app{
		Templates: template.Must(template.ParseFS(TmplFiles, "ui/tmpl/*.html")),
	}

	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/healthcheck"))
	r.Use(middleware.Timeout(3 * time.Second))

	r.Get("/", app.index)

	// Serve remaining static files
	sub, err := fs.Sub(StaticFiles, "ui")
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.FS(sub))
	//r.Mount("/static", http.StripPrefix("ui/static", fs))
	r.Mount("/static", fs)

	// HTTPS server struct utilizing application
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Do the thing
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

type app struct {
	Templates *template.Template
}

func (app *app) index(w http.ResponseWriter, r *http.Request) {
	app.Templates.ExecuteTemplate(w, "index.html", nil)
}
