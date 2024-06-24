package bill

import "fmt"

type Country int

const (
	SERBIA Country = iota // int 0
	TURKEY
	RUSSIA
)

// must align with the Country enum
var countryToString = []string{"serbia", "turkey", "russia"}

func (c Country) String() string {
	return countryToString[c]
}

func StringToCountry(countryString string) (Country, error) {
	switch countryString {
	case "serbia":
		return SERBIA, nil
	case "turkey":
		return TURKEY, nil
	case "russia":
		return RUSSIA, nil
	default:
		return -1, fmt.Errorf("Country %s not found", countryString)
	}
}

func GetCountryList() []string {
	countryList := append([]string{}, countryToString...)
	return countryList
}
