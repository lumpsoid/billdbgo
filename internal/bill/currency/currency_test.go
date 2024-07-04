package currency

import (
	"testing"
)

func TestCurrencyToString(t *testing.T) {
	currency := RSD
	if currency.String() != "rsd" {
		t.Errorf("Currency as string: `%s`, expected `rsd`", currency)
	}
	currency = RUB
	if currency.String() != "rub" {
		t.Errorf("Currency as string: `%s`, expected `rub`", currency)
	}
}

func TestParse(t *testing.T) {
	currencyString := "usd"
	currency, err := Parse(currencyString)
	if err != nil {
		t.Error("Error parsing currency string:", err)
	}
	if currency.String() != currencyString {
		t.Errorf("Currency: `%s`, expected `%s`", currency, currencyString)
	}
	currencyString = "rub"
	currency, err = Parse(currencyString)
	if err != nil {
		t.Error("Error parsing currency string:", err)
	}
	if currency.String() != currencyString {
		t.Errorf("Currency: `%s`, expected `%s`", currency, currencyString)
	}
}

func TestAvailable(t *testing.T) {
	currencies := Available()
	if len(currencies) != len(currencyToString) {
		t.Errorf("Currencies length: %d, expected %d",
			len(currencies),
			len(currencyToString),
		)
	}
	if currencies[3] != currencyToString[3] {
		t.Errorf(
			"Currency: `%s`, expected `%s`",
			currencies[0],
			currencyToString[0],
		)
	}
}
