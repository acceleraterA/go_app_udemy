package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/acceleraterA/go_app_udemy/internal/config"
	"github.com/acceleraterA/go_app_udemy/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a

}
func AddDefaultData(td *models.TemplateData, req *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(req)
	return td
}
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData, r *http.Request) {

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
		log.Fatal("Could not get template from template cache")
	}
	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	_ = parsedTemplate.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}
func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache:=make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	//get all of the files named *.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
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
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
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

func RenderTemplate(w http.ResponseWriter, t string) {
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
