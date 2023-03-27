package main

import (
	"fmt"
	"testing"

	"github.com/acceleraterA/go_app_udemy/internal/config"
	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//
	default:
		t.Errorf(fmt.Sprintf("type is not *chi.Mux, but is %T", v))
	}
}
