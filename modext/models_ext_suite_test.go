package modext_test

import (
	"database/sql"
	"regexp"

	"github.com/davecgh/go-spew/spew"
	"github.com/harrisbaird/dailyteedeals/config"
	"github.com/harrisbaird/dailyteedeals/models"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
)

func RunInTestTransaction(fn func(boil.Executor)) {
	db := newTestDB()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	fn(tx)
}

func TableCountDiff(db boil.Executor, table string, testFn func()) int {
	countBefore := getTableCount(db, table)
	testFn()
	return getTableCount(db, table) - countBefore
}

func CreateArtistFixtures(db boil.Executor) models.ArtistSlice {
	artists := models.ArtistSlice{
		&models.Artist{ID: 1, Name: "theduc", Urls: []string{"http://teepublic.com/user/theduc", "http://neatoshop.com/artist/Theduc"}},
	}
	for _, artist := range artists {
		err := artist.Upsert(db, false, []string{"id"}, []string{})
		if err != nil {
			panic(err)
		}
	}
	return artists
}

func CreateDesignFixtures(db boil.Executor) models.DesignSlice {
	artists := CreateArtistFixtures(db)
	designs := models.DesignSlice{
		&models.Design{ID: 1, ArtistID: artists[0].ID, Name: "Summer is here"},
		&models.Design{ID: 2, ArtistID: artists[0].ID, Name: "Training Corps"},
		&models.Design{ID: 3, ArtistID: artists[0].ID, Name: "Ninja"},
		&models.Design{ID: 4, ArtistID: artists[0].ID, Name: "Survey Corps"},
	}
	for _, design := range designs {
		err := design.Upsert(db, false, []string{"id"}, []string{})
		if err != nil {
			panic(err)
		}
	}
	return designs
}

func CreateProductFixtures(db boil.Executor) models.ProductSlice {
	designs := CreateDesignFixtures(db)
	sites := CreateSiteFixtures(db)
	products := models.ProductSlice{
		&models.Product{ID: 1, DesignID: designs[0].ID, SiteID: sites[0].ID, Active: true, Deal: true},
		&models.Product{ID: 2, DesignID: designs[1].ID, SiteID: sites[0].ID, Active: true, Deal: false},
		&models.Product{ID: 3, DesignID: designs[2].ID, SiteID: sites[1].ID, Active: true, Deal: true},
		&models.Product{ID: 4, DesignID: designs[3].ID, SiteID: sites[1].ID, Active: true, Deal: false},
	}
	for _, product := range products {
		err := product.Upsert(db, false, []string{"id"}, []string{"id", "design_id", "site_id", "slug", "active", "deal"})
		if err != nil {
			spew.Dump(product)
			panic(err)
		}
	}
	return products
}

func CreateSiteFixtures(db boil.Executor) models.SiteSlice {
	sites := models.SiteSlice{
		&models.Site{ID: 1, Name: "Qwertee", DomainName: "qwertee.com"},
		&models.Site{ID: 2, Name: "Teefury", DomainName: "teefury.com"},
	}
	for _, site := range sites {
		err := site.Upsert(db, false, []string{"id"}, []string{})
		if err != nil {
			panic(err)
		}
	}
	return sites
}

func CreateUserFixtures(db boil.Executor) models.UserSlice {
	users := models.UserSlice{
		&models.User{ID: 1, Email: "active_api_user@test.com", APIAccess: true, APIToken: "active_api_user"},
		&models.User{ID: 2, Email: "inactive_api_user@test.com", APIAccess: false, APIToken: "inactive_api_user"},
	}
	for _, user := range users {
		err := user.Upsert(db, false, []string{"id"}, []string{})
		if err != nil {
			panic(err)
		}
	}
	return users
}

var ValidSlug = regexp.MustCompile(`^\d{5}-[a-z0-9-]+`)

func newTestDB() *sql.DB {
	db, err := sql.Open("postgres", config.PostgresTestConnectionString())
	if err != nil {
		panic(err)
	}

	return db
}

func getTableCount(db boil.Executor, table string) int {
	type Count struct {
		Count int `boil:"count"`
	}
	var count Count
	err := models.NewQuery(db, qm.Select("count(*) as count"), qm.From(table)).Bind(&count)
	if err != nil {
		panic(err)
	}

	return count.Count
}
