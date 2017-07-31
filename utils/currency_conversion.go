package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"sort"
)

const DefaultCurrency = "USD"

type ApproximatePrice struct {
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Approximate bool   `json:"approximate"`
	Formatted   string `json:"formatted"`
}

type Currency struct {
	Symbol string
	Rates  map[string]float64 `json:"rates"`
}

var CurrentRates = map[string]*Currency{
	"USD": &Currency{Symbol: "$"},
	"GBP": &Currency{Symbol: "£"},
	"EUR": &Currency{Symbol: "€"},
}

func UpdateRates() error {
	wantedCurrencies := strings.Join(getCurrencies(), ",")

	for key, currency := range CurrentRates {
		resp, err := HTTPGetBytes(fmt.Sprintf("https://api.fixer.io/latest?base=%s&symbols=%s", key, wantedCurrencies))
		if err != nil {
			return err
		}

		if err := json.Unmarshal(resp, &currency); err != nil {
			return err
		}
	}
	return nil
}

func ConvertPrices(rawPrices map[string]string) map[string]*ApproximatePrice {
	keys := []string{}
	output := make(map[string]*ApproximatePrice)

	for k, v := range rawPrices {
		number, _ := strconv.Atoi(v)
		key := strings.ToUpper(k)
		output[key] = &ApproximatePrice{Amount: number, Currency: key, Formatted: format(number, key)}
		keys = append(keys, key)
	}

	wanted := StringSlicesDiff(keys, getCurrencies())

	// Add any missing currencies and mark as approximate
	for _, to := range wanted {
		from := keys[0]
		amount := output[from].Amount

		rate := CurrentRates[from].Rates[to]
		convertedAmount := int(float64(amount) * rate)
		output[to] = &ApproximatePrice{Amount: convertedAmount, Currency: to, Formatted: format(convertedAmount, to), Approximate: true}
	}

	return output
}

func getCurrencies() []string {
	currencies := []string{}
	for key := range CurrentRates {
		currencies = append(currencies, key)
	}
	sort.Strings(currencies)
	return currencies
}

func format(amount int, currency string) string {
	return fmt.Sprintf("%s%.0f", CurrentRates[currency].Symbol, float64(amount)/100)
}
