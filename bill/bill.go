package bill

import (
	"time"
)

type Currency int

const (
	EUR Currency = iota
	RSD
	TRY
	RUB
)

type Country int

const (
	SERBIA Country = iota
	TURKEY
	RUSSIA
)

type Bill struct {
	Name         string
	Date         time.Time
	Price        int
	Currency     Currency
	ExchangeRate int
	Country      Country
	Items        []Item
	Tag          Tag
}
