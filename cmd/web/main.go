package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/acceleraterA/go_app_udemy/internal/config"
	"github.com/acceleraterA/go_app_udemy/internal/driver"
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
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	//close the database when the main(app) is stopped running
	defer db.SQL.Close()
	fmt.Printf(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
func run() (*driver.DB, error) {
	// (register the reservation object to session) what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
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
	//connect to database
	log.Println("connecting to database...")
	//password will be updated later
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=")
	if err != nil {
		log.Fatal("cannot connect to db, dying...")
	}
	log.Println("Connected to database!")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false
	// give render access to app
	render.NewRenderer(&app)
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)
	helpers.NewHelper(&app)
	return db, nil
}
