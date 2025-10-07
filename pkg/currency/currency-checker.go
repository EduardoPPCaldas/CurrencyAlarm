package currency

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type CurrencyChecker interface {
	CheckCurrency(ownedCurrency string, exchangedCurrency string) (float64, error)
}

type awesomeCurrencyChecker struct {
}

func NewCurrencyChecker() CurrencyChecker {
	return &awesomeCurrencyChecker{}
}

func (cc awesomeCurrencyChecker) CheckCurrency(ownedCurrency string, convertedCurrency string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/%s-%s", convertedCurrency, ownedCurrency))
	if err != nil {
		fmt.Println("Error fetching:", err)
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]struct {
		Bid string `json:"bid"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding:", err)
		return 0, nil
	}

	conversion := fmt.Sprintf("%s%s", convertedCurrency, ownedCurrency)
	bidStr := result[conversion].Bid
	bid, err := strconv.ParseFloat(bidStr, 64)
	if err != nil {
		fmt.Println("Error parsing bid:", err)
		return 0, nil
	}

	return bid, nil
}
