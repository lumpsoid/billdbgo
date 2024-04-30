package bill

import "fmt"

type Country int

const (
	SERBIA Country = iota
	TURKEY
	RUSSIA
)

// Map string values to enum values
var countryMap = map[string]Country{
	"serbia": SERBIA,
	"turkiye": TURKEY,
	"russia": RUSSIA,
}

var countryToString = []string{"serbia", "turkey", "russia"}

func StringToCountry(s string) (Country, error) {
	if country, ok := countryMap[s]; ok {
		return country, nil
	}
	return -1, fmt.Errorf("Country %s not found", s)
}
