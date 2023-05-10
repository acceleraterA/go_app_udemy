package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/acceleraterA/go_app_udemy/internal/config"
	"github.com/acceleraterA/go_app_udemy/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"

//var functions = template.FuncMap{}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a

}
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	//flush error and warning will be automatically populated when we rendering the templates
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}
func Template(w http.ResponseWriter, tmpl string, td *models.TemplateData, r *http.Request) error {

	var tc map[string]*template.Template
	if app.UseCache {
		//get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get requested template from cache
	parsedTemplate, ok := tc[tmpl]
	if !ok {
		//log.Fatal("Could not get template from template cache")
		return errors.New("Can't get template from cache")
	}
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = parsedTemplate.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}
	return nil
}
func CreateTemplateCache() (map[string]*template.Template, error) {
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
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

/*
// package level variable tc (parsed template cache)
var tc = make(map[string]*template.Template)

func Template(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check to see if we already have the template in cache
	_, inMap := tc[t]
	if !inMap {
		//need to create the template
		log.Println("creating template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		//have the template in cache
		log.Println("using cached template")
	}
	tmpl = tc[t]
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}
}


func createTemplateCache(t string) error {
	//a function to parse template and save in cache
	templates := []string{"./templates/" + t, "./templates/base.layout.tmpl"}
	//parse the template
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	//add template to cache
	tc[t] = tmpl
	return nil
}
*/
