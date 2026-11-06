package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

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

	r := newRouter(app)

	// HTTP server; TLS is terminated at the reverse proxy (Cloudflare) in front of this service
	srv := &http.Server{
		Addr:    ":8080", // TODO(config): drive from environment configuration, see TODO.md
		Handler: r,
	}

	// Start the HTTP server
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

// newRouter wires up the chi router with all middleware, routes and handlers.
func newRouter(app *app) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/healthcheck"))

	r.Get("/", app.index)

	// Serve remaining static files. fs.Sub strips the "ui" embedded prefix,
	// so no additional StripPrefix is needed.
	sub, err := fs.Sub(StaticFiles, "ui")
	if err != nil {
		panic(err)
	}
	r.Mount("/static", http.FileServer(http.FS(sub)))

	// Load static root directory files e.g. favicon
	// Note: This uses the RootFiles embedded FS for favicon and other root assets
	app.loadStaticRoutes(r, RootFiles, "ui/root")

	return r
}

type app struct {
	Templates *template.Template
}

func (app *app) index(w http.ResponseWriter, r *http.Request) {
	if err := app.Templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("error executing index template: %v", err)
	}
}

// rootAssetExts is the allowlist of file extensions loadStaticRoutes serves
// from ui/root/. Keep it in sync with the web-compliance assets the root
// directory is expected to hold over time: PWA/android/apple icons, web app
// manifest, résumé, and the like.
var rootAssetExts = map[string]bool{
	".ico":         true,
	".png":         true,
	".svg":         true,
	".webmanifest": true,
	".json":        true,
	".pdf":         true,
}

func (app *app) loadStaticRoutes(r *chi.Mux, fileSys fs.FS, path string) {
	rootEntries, err := fs.ReadDir(fileSys, path)
	if err != nil {
		log.Printf("Error reading root directory %s: %v", path, err)
		return
	}

	for _, entry := range rootEntries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// File type validation (allowlist)
		ext := strings.ToLower(filepath.Ext(name))
		if !rootAssetExts[ext] {
			log.Printf("Disallowed file type: %s", name)
			continue
		}

		// Simple path traversal check
		if strings.Contains(name, "..") || strings.HasPrefix(name, "/") {
			log.Printf("Path traversal attempt detected: %s", name)
			continue
		}

		safePath := path + "/" + name
		r.Get("/"+name, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, fileSys, safePath)
		})
	}
}
