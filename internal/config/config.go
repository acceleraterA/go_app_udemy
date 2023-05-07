package config

import (
	"html/template"
	"log"

	"github.com/acceleraterA/go_app_udemy/internal/models"
	scs "github.com/alexedwards/scs/v2"
)

//doesn't import any app to avoid import cycle

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
