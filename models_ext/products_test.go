package models_ext_test

import (
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	. "github.com/harrisbaird/dailyteedeals/models_ext"
	"github.com/nbio/st"
)

// func TestX(t *testing.T) {
// 	db, mock := newSQLMock()
// 	defer db.Close()

// 	mock.ExpectExec("SELECT \"products\".* FROM \"products\" INNER JOIN sites on products.site_id = sites.id WHERE (products.active = true) AND (products.deal = true) ORDER BY sites.display_order ASC, products.site_id ASC, products.last_chance ASC, products.slug ASC;")

// 	spew.Dump(ActiveDeals(db))

// }

func TestMarkProductsInactive(t *testing.T) {
	type args struct {
		siteID int
		deal   bool
	}

	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"deal", args{1, true}, false},
		{"non-deal", args{1, false}, false},
	}

	db, mock := newSQLMock()
	defer db.Close()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("UPDATE \"products\"").
				WithArgs(false, tt.args.siteID, tt.args.deal).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := MarkProductsInactive(db, tt.args.siteID, tt.args.deal)
			st.Expect(t, err != nil, tt.wantErr)
			st.Expect(t, mock.ExpectationsWereMet() != nil, false)
		})
	}
}
