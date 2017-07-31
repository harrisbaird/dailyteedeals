package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/orm"
	"github.com/harrisbaird/dailyteedeals/models"
	. "github.com/harrisbaird/dailyteedeals/server"
	"github.com/harrisbaird/dailyteedeals/utils"
	"github.com/nbio/st"
)

func TestProductRedirectRouter(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		routerPath string
		wantURL    string
		wantErr    bool
	}{
		{"Valid - Non affiliate", 302, "/non_affiliate_link", "http://test.com", false},
		{"Valid - affiliate", 302, "/affiliate_link", "http://affiliatesite.com?url=http%3A%2F%2Ftest.com", false},
		{"Invalid", 404, "/invalid", "", false},
	}

	models.RunInTestTransaction(false, func(db orm.DB) {
		models.ImportProductFixtures(db)

		server := httptest.NewServer(ProductRedirectRouter(db, gin.New()))
		defer server.Close()

		client := http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		for _, tt := range testCases {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := client.Get(server.URL + tt.routerPath)
				st.Assert(t, err, nil)
				defer resp.Body.Close()

				st.Assert(t, resp.StatusCode, tt.statusCode)

				_, err = ioutil.ReadAll(resp.Body)
				st.Assert(t, err, nil)
				st.Assert(t, resp.Header.Get("Location"), tt.wantURL)
			})
		}
	})

}

func TestApi(t *testing.T) {
	testCases := []struct {
		name        string
		statusCode  int
		routerPath  string
		fixturePath string
	}{
		{"V1 Products", 200, "/v1/products.json?key=active_api_user", "v1/products.json"},

		{"V2 Deals", 200, "/v2/deals?key=active_api_user", "v2/deals.json"},
		{"V2 Site", 200, "/v2/sites/qwertee?key=active_api_user", "v2/site.json"},
		{"V2 Artist", 200, "/v2/artists/55555-theduc?key=active_api_user", "v2/artist.json"},
		{"V2 Design", 200, "/v2/designs/55555-summer-is-here?key=active_api_user", "v2/design.json"},

		{"V2 Site - Invalid", 404, "/v2/sites/invalid?key=active_api_user", "v2/site_invalid.json"},
		{"V2 Artist - Invalid", 404, "/v2/artists/invalid?key=active_api_user", "v2/artist_invalid.json"},
		{"V2 Design - Invalid", 404, "/v2/designs/invalid?key=active_api_user", "v2/design_invalid.json"},

		{"Non api user", 401, "/v1/products.json?key=inactive_api_user", "non_api_user.json"},
		{"Invalid Api Key", 401, "/v1/products.json?key=invalid", "invalid_key.json"},
		{"Missing Api Key", 401, "/v1/products.json", "missing_key.json"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			models.RunInTestTransaction(false, func(db orm.DB) {
				models.ImportUserFixtures(db)
				models.ImportProductFixtures(db)

				server := httptest.NewServer(ApiRouter(db, gin.New()))
				defer server.Close()
				resp, err := http.Get(server.URL + tt.routerPath)
				st.Assert(t, err, nil)
				defer resp.Body.Close()

				have, err := ioutil.ReadAll(resp.Body)

				st.Expect(t, resp.StatusCode, tt.statusCode)
				st.Expect(t, resp.Header.Get("Content-Type"), "application/json; charset=utf-8")

				st.Assert(t, err, nil)
				want, err := ioutil.ReadFile(utils.ProjectRootPath("server", "testdata", tt.fixturePath))
				st.Assert(t, err, nil)

				st.Expect(t, have, want)
			})
		})
	}
}
