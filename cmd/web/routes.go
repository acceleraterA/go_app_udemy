package main

import (
	"net/http"

	"github.com/acceleraterA/go_app_udemy/pkg/config"
	"github.com/acceleraterA/go_app_udemy/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(NoSurf)
	r.Use(SessionLoad)
	r.Get("/", handlers.Repo.Home)
	r.Get("/about", handlers.Repo.About)
	fileServer := http.FileServer(http.Dir("./static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return r
}
