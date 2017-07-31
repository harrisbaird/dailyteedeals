package utils_test

import (
	"testing"

	gock "gopkg.in/h2non/gock.v1"

	"fmt"
	"strings"

	. "github.com/harrisbaird/dailyteedeals/utils"
	"github.com/nbio/st"
)

func TestConvertPrices(t *testing.T) {
	defer gock.Off()

	for _, currency := range []string{"usd", "gbp", "eur"} {
		gock.New("https://api.fixer.io" + fmt.Sprintf("/latest?base=%s&symbols=EUR,GBP,USD", strings.ToUpper(currency))).
			Reply(200).
			File(ProjectRootPath("utils", "testdata", "currency", currency+".json"))
	}

	err := UpdateRates()
	st.Assert(t, err, nil)

	prices := ConvertPrices(map[string]string{"USD": "1200"})
	st.Expect(t, prices, map[string]ApproximatePrice{
		"USD": ApproximatePrice{Amount: 1200, Currency: "USD", Formatted: "$12", Approximate: false},
		"EUR": ApproximatePrice{Amount: 1026, Currency: "EUR", Formatted: "€10", Approximate: true},
		"GBP": ApproximatePrice{Amount: 913, Currency: "GBP", Formatted: "£9", Approximate: true},
	})
	st.Expect(t, gock.HasUnmatchedRequest(), false)
}
