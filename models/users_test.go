package models_test

import (
	"testing"

	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/database"
	. "github.com/harrisbaird/dailyteedeals/models"
	"github.com/nbio/st"
)

func TestValidAPIUser(t *testing.T) {
	testCases := []struct {
		name      string
		token     string
		wantValid bool
	}{
		{"active", "active_api_user", true},
		{"inactive", "inactive_api_user", false},
		{"invalid", "invalid", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			database.RunInTestTransaction(false, func(db orm.DB) {
				ImportUserFixtures(db)
				st.Expect(t, ValidAPIUser(db, tt.token), tt.wantValid)
			})
		})
	}
}
