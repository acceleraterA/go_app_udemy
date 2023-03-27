package render

import (
	"net/http"
	"testing"

	"github.com/acceleraterA/go_app_udemy/internal/models"
)

//buid request with session data

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("failed")
	}
}
func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter
	err = RenderTemplate(&ww, "home.page.tmpl", &models.TemplateData{}, r)
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = RenderTemplate(&ww, "home1.page.tmpl", &models.TemplateData{}, r)
	if err == nil {
		t.Error("rendered template that does not exist")
	}

}
func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	//add session to r.context
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	//update request
	r = r.WithContext(ctx)
	return r, nil
}

func TestNewTemplate(t *testing.T) {
	NewTemplates(app)

}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
