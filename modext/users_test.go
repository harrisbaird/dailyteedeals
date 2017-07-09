package modext_test

import (
	"testing"

	. "github.com/harrisbaird/dailyteedeals/modext"
	"github.com/nbio/st"
	"github.com/vattle/sqlboiler/boil"
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
			RunInTestTransaction(func(db boil.Executor) {
				CreateUserFixtures(db)
				st.Expect(t, ValidAPIUser(db, tt.token), tt.wantValid)
			})
		})
	}
}
