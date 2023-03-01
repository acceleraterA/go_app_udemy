package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/acceleraterA/go_app_udemy/pkg/config"
	"github.com/acceleraterA/go_app_udemy/pkg/handlers"
	"github.com/acceleraterA/go_app_udemy/pkg/render"
)

const portNumber = ":8080"

func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false
	render.NewTemplates(&app)
	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
