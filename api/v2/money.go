package v2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/harrisbaird/dailyteedeals/utils"
	"github.com/vattle/sqlboiler/types"
)

var wantedCurrencies = []string{"USD", "GBP", "EUR"}

type ApproximatePrice struct {
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Approximate bool   `json:"approximate"`
}

// TODO: Automatically update
var currencyRates = map[string]float64{
	"USD2GBP": 0.78,
	"USD2EUR": 0.9,
	"GBP2USD": 1.27,
	"GBP2EUR": 1.14,
	"EUR2USD": 1.12,
	"EUR2GBP": 0.88,
}

var CurrencySymbols = map[string]string{
	"USD": "$",
	"GBP": "£",
	"EUR": "€",
}

func ConvertPrices(rawPrices types.HStore) map[string]ApproximatePrice {
	keys := []string{}
	output := make(map[string]ApproximatePrice)

	for k, v := range rawPrices {
		number, _ := strconv.Atoi(v.String)
		key := strings.ToUpper(k)
		output[key] = ApproximatePrice{Amount: number, Currency: key}
		keys = append(keys, key)
	}

	wanted := utils.StringSlicesDiff(keys, wantedCurrencies)

	// Add any missing currencies and mark as approximate
	for _, key := range wanted {
		from := keys[0]
		amount := output[from].Amount

		output[key] = ApproximatePrice{Amount: convert(amount, from, key), Currency: key, Approximate: true}
	}

	return output
}

func getKeys(mymap map[string]ApproximatePrice) []string {
	keys := []string{}

	for k := range mymap {
		keys = append(keys, k)
	}

	return keys
}

func convert(amount int, from string, to string) int {
	key := fmt.Sprintf("%s2%s", from, to)
	rate := currencyRates[key]

	return int(float64(amount) * rate)
}
