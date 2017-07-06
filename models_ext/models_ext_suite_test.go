package models_ext_test

import (
	"database/sql"
	"regexp"

	"log"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func newSQLMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

var validSlug = regexp.MustCompile(`^\d{5}-[a-z0-9-]+`)

func validateSlug(slug string) bool {
	return validSlug.MatchString(slug)
}
