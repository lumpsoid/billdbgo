package bill

import "fmt"

type Currency int

const (
	EUR Currency = iota
	RSD
	TRY
	RUB
)

// Map string values to enum values
var currencyMap = map[string]Currency{
	"eur": EUR,
	"rsd": RSD,
	"try": TRY,
}

var currencyToString = []string{"eur", "rsd", "try", "rub"}

func StringToCurrency(s string) (Currency, error) {
	if currency, ok := currencyMap[s]; ok {
		return currency, nil
	}
	return -1, fmt.Errorf("Currency %s not found", s)
}
