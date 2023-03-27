package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/acceleraterA/go_app_udemy/internal/config"
	"github.com/acceleraterA/go_app_udemy/internal/handlers"
	"github.com/acceleraterA/go_app_udemy/internal/helpers"
	"github.com/acceleraterA/go_app_udemy/internal/models"
	"github.com/acceleraterA/go_app_udemy/internal/render"
	scs "github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog log.Logger
var ErrorLog log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
func run() error {
	// (register the reservation object to session) what am I going to put in the session
	gob.Register(models.Reservation{})
	//change this to true when in production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	ErrorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = ErrorLog

	// Initialize a new session manager and configure the session lifetime.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	// store the session to config app.Session
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
		return err
	}
	app.TemplateCache = tc
	app.UseCache = false
	// give render access to app
	render.NewTemplates(&app)
	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)
	helpers.NewHelper(&app)
	return nil
}
