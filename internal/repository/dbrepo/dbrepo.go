package dbrepo

import (
	"database/sql"

	"github.com/acceleraterA/go_app_udemy/internal/config"
	"github.com/acceleraterA/go_app_udemy/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
