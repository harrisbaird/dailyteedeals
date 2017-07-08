package models_ext_test

import (
	"testing"

	"github.com/harrisbaird/dailyteedeals/models"
	. "github.com/harrisbaird/dailyteedeals/models_ext"
	"github.com/nbio/st"
	"github.com/vattle/sqlboiler/boil"
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
			RunInTestTransaction(func(db boil.Executor) {
				CreateDesignFixtures(db)
				var design *models.Design
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

func TestDesignHooks(t *testing.T) {
	RunInTestTransaction(func(db boil.Executor) {
		CreateArtistFixtures(db)
		design := models.Design{
			ArtistID:     1,
			Name:         "   test design  ",
			Tags:         []string{"tags 1"},
			CategoryTags: []string{"category tags 1"}}
		err := design.Insert(db)

		st.Expect(t, err, nil)
		st.Expect(t, ValidSlug.MatchString(design.Slug), true)
	})
}
