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

//go:embed all:ui/root
var RootFiles embed.FS

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
	//r.Get("/favicon.ico", app.favicon)

	// Serve remaining static files
	sub, err := fs.Sub(StaticFiles, "ui")
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.FS(sub))
	//r.Mount("/static", http.StripPrefix("ui/static", fs))
	r.Mount("/static", fs)

	// Load static root directory files e.g. favicon
	app.loadStaticRoutes(r, RootFiles, "ui/root")

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

func (app *app) favicon(w http.ResponseWriter, r *http.Request) {
	// StaticFiles is rooted at "ui/static/..." because of how go:embed
	// keeps the path it was given — this must match that root exactly.
	http.ServeFileFS(w, r, StaticFiles, "ui/static/media/favicon.ico")
}

func (app *app) loadStaticRoutes(r *chi.Mux, fileSys fs.FS, path string) {
	rootEntries, err := fs.ReadDir(fileSys, path)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range rootEntries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		r.Get("/"+name, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, fileSys, path+"/"+name)
		})
	}
}
