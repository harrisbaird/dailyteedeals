package models_test

import (
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/database"
	. "github.com/harrisbaird/dailyteedeals/models"
	"github.com/nbio/st"
)

func TestFindOrCreateDesign(t *testing.T) {
	testCases := []struct {
		name               string
		designName         string
		designsCountChange int
	}{
		{"New", "Susuwatari family", 1},
		{"Existing", "Summer is here", 0},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			database.RunInTestTransaction(false, func(db orm.DB) {
				ImportDesignFixtures(db)
				var design *Design
				var err error

				count := TableCountDiff(db, "designs", func() {
					design, err = FindOrCreateDesign(db, 1, tt.designName)
				})

				st.Expect(t, err, nil)
				st.Reject(t, design, nil)
				st.Expect(t, count, tt.designsCountChange)
				st.Expect(t, ValidSlug.MatchString(design.Slug), true)
			})
		})
	}
}

func TestFindDesignBySlug(t *testing.T) {
	testCases := []struct {
		name    string
		slug    string
		wantID  int
		wantErr bool
	}{
		{"Found", "55555-summer-is-here", 1, false},
		{"Missing", "missing", 0, true},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			database.RunInTestTransaction(false, func(db orm.DB) {
				ImportDesignFixtures(db)
				design, err := FindDesignBySlug(db, tt.slug)
				st.Expect(t, design.ID, tt.wantID)
				st.Expect(t, err != nil, tt.wantErr)
			})
		})
	}
}

// func TestDesignHooks(t *testing.T) {
// 	database.RunInTestTransaction(false, func(db boil.Executor) {
// 		CreateArtistFixtures(db)
// 		design := models.Design{
// 			ArtistID:     1,
// 			Name:         "   test design  ",
// 			Tags:         []string{"tags 1"},
// 			CategoryTags: []string{"category tags 1"}}
// 		err := design.Insert(db)

// 		st.Expect(t, err, nil)
// 		st.Expect(t, ValidSlug.MatchString(design.Slug), true)
// 	})
// }
