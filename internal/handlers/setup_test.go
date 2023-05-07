package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/acceleraterA/go_app_udemy/internal/config"
	"github.com/acceleraterA/go_app_udemy/internal/models"
	"github.com/acceleraterA/go_app_udemy/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func TestMain(m *testing.M) {
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
	//connect to database
	log.Println("connecting to database...")
	/*db, err := driver.ConnectSQL("host=localhost port=5432 dbname=test_connect user=postgres password=Bastille8877,,")
	  if err != nil {
	  	log.Fatal("cannot connect to db, dying...")
	  }
	  log.Println("Connected to database!")
	*/
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(app.MailChan)
	fmt.Println("starting mail listener...")
	//duplicate the function in actual application
	listenForMail()
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
		//return err
	}
	app.TemplateCache = tc
	app.UseCache = true
	// give render access to app

	repo := NewTestRepo(&app)
	NewHandler(repo)
	render.NewRenderer(&app)

	os.Exit(m.Run())

}
func getRoutes() http.Handler {
	//copy from routes.go
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	//r.Use(NoSurf)
	r.Use(SessionLoad)

	r.Get("/", Repo.Home)
	r.Get("/about", Repo.About)
	r.Get("/generals-quarters", Repo.Generals)
	r.Get("/majors-suite", Repo.Majors)

	r.Get("/search-availability", Repo.Availability)
	r.Post("/search-availability", Repo.PostAvailability)
	r.Post("/search-availability-json", Repo.AvailabilityJSON)

	r.Get("/contact", Repo.Contact)
	r.Get("/make-reservation", Repo.Reservation)
	r.Post("/make-reservation", Repo.PostReservation)
	r.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return r
}

// Nosurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	//myCache:=make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	//get all of the files named *.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}
	//ranging through all files ending with *.page.tmpl
	for _, page := range pages {
		//get the name of page
		name := filepath.Base(page)
		//create new template from page named name
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

// skip actual sending mail
func listenForMail() {
	go func() {
		//infinate for loop
		for {
			_ = <-app.MailChan

		}
	}()
}
