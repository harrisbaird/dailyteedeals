package migrations

import (
	"log"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/harrisbaird/dailyteedeals/database"
)

func Run() {
	db := database.Connect()
	defer db.Close()
	log.Println("Migrating")
	runMigrations(db)
}

func RunTest() {
	log.Println("Migrating Test")
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
