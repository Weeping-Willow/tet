package storage

import (
	"embed"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Migrate(db *sqlx.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		return errors.Wrap(err, "set goose dialect")
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return errors.Wrap(err, "migrate database")
	}

	return nil
}
