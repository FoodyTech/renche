package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type App struct {
	cache Cache
}

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}) error
}

func New() *App {
	a := &App{}
	return a
}

func (app *App) Run() error {
	log.Println("run...")

	srv := &http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      app.initRoutes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}
	return srv.ListenAndServe()
}

func (app *App) initRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/status/ping"))
	r.Use(middleware.RealIP)
	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(1 * time.Minute))

	r.Route("/", func(r chi.Router) {
		r.HandleFunc("/render/{format}", app.renderHandler)

		r.Route("/pages", func(r chi.Router) {
			r.HandleFunc("/{id}/status", app.pageStatusHandler)
			r.HandleFunc("/{id}/result/{format}", app.pageResultHandler)
		})

		r.Route("/sitemap", func(r chi.Router) {
			r.HandleFunc("/xml/url", app.sitemapHandler)
			r.HandleFunc("/xml/file", app.sitemapFileHandler)
		})
	})
	return r
}
