package bill

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
