package country

import (
	"fmt"
	"testing"
)

func TestCountryToString(t *testing.T) {
	fmt.Println("TestCountryToString")
	country := SERBIA
	if country.String() != "serbia" {
		t.Errorf("Currency as string: `%s`, expected `rsd`", country)
	}
	country = TURKEY
	if country.String() != "turkey" {
		t.Errorf("Currency as string: `%s`, expected `rub`", country)
	}
}
