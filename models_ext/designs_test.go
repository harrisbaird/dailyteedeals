package models_ext_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/harrisbaird/dailyteedeals/models"
	. "github.com/harrisbaird/dailyteedeals/models_ext"
	"github.com/nbio/st"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFindOrCreateDesign(t *testing.T) {
	db, mock := newSQLMock()
	defer db.Close()

	testCases := []struct {
		name      string
		createNew bool
	}{
		{"New", true},
		{"Existing", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			rows := sqlmock.NewRows([]string{"id", "name", "slug", "tags"})

			if tt.createNew {
				mock.ExpectQuery("SELECT \\* FROM \"designs\"").WillReturnRows(rows)
				mock.ExpectQuery("INSERT INTO \"designs\"").
					WithArgs(1, "Design Name", sqlmock.AnyArg(), "{\"designname\"}", "{}").
					WillReturnRows(sqlmock.NewRows([]string{"id", "description", "mature", "active_products_count"}).AddRow(1, "", false, 0))
			} else {
				mock.ExpectQuery("SELECT \\* FROM \"designs\"").WillReturnRows(rows.AddRow(1, "Design Name", "45478-designname", "{\"designname\"}"))
			}

			design, err := FindOrCreateDesign(db, 1, "Design Name")
			spew.Dump(design, err)

			st.Expect(t, err, nil)
			st.Expect(t, validateSlug(design.Slug), true)
			st.Expect(t, mock.ExpectationsWereMet() != nil, false)
		})
	}
}

func TestDesignHooks(t *testing.T) {
	db, mock := newSQLMock()
	defer db.Close()

	design := models.Design{
		ArtistID:     1,
		Name:         "   test design  ",
		Tags:         []string{"tags 1"},
		CategoryTags: []string{"category tags 1"}}

	mock.ExpectQuery("INSERT INTO \"designs\"").
		WithArgs(1, "Test Design", sqlmock.AnyArg(), "{\"tags1\",\"testdesign\"}", "{\"categorytags1\"}").
		WillReturnRows(sqlmock.NewRows([]string{"id", "description", "mature", "active_products_count"}).AddRow(1, "", false, 0))

	err := design.Insert(db)

	st.Expect(t, err, nil)
	st.Expect(t, validateSlug(design.Slug), true)
	st.Expect(t, mock.ExpectationsWereMet(), nil)
}
