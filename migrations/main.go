package migrations

import (
	"log"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/harrisbaird/dailyteedeals/database"
)

func Run() {
	log.Println("Migrating database")
	db := database.Connect()
	defer db.Close()
	runMigrations(db)
}

func RunTest() {
	log.Println("Migrating test database")
	db := database.ConnectTest()
	defer db.Close()
	runMigrations(db)
}

func runMigrations(db *pg.DB) {
	migrations.Run(db, "init")
	_, _, err := migrations.Run(db, "up")
	if err != nil {
		panic(err)
	}
}
