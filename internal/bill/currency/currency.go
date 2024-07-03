package currency

import "fmt"

type Currency int

const (
	EUR Currency = iota // int 0
	RSD
	TRY
	RUB
	USD
)

// must be aligned with the Currency enum
var currencyToString = []string{
	"eur", "rsd", "try", "rub", "usd",
}

func (c Currency) String() string {
	return currencyToString[c]
}

func Available() []string {
	currencyList := append([]string{}, currencyToString...)
	return currencyList
}

func Parse(currencyStr string) (Currency, error) {
	switch currencyStr {
	case "usd":
		return USD, nil
	case "eur":
		return EUR, nil
	case "rsd":
		return RSD, nil
	case "try":
		return TRY, nil
	case "rub":
		return RUB, nil
	default:
		return -1, fmt.Errorf("Currency %s not found", currencyStr)
	}
}
