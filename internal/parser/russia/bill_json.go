package russia

import (
	B "billdb/internal/bill"
	"time"
)

type BillJson struct {
	Code    int
	First   int
	Data    Data
	Request Request
}

type Data struct {
	Json Json
	Html string
}

type Json struct {
	Items       []Items
	DateTime    string
	KktRegId    string
	RetailPlace string
	TotalSum    int
}

type Items struct {
	Nds                  int
	Sum                  int
	Name                 string
	Price                int
	Quantity             float64
	PaymentType          int
	ProductType          int
	ItemsQuantityMeasure int
}

type Request struct {
	Qrurl  string
	Qrfile string
	Qrraw  string
	Manual Manual
}

type Manual struct {
	Fn        string
	Fd        string
	Fp        string
	CheckTime string
	Type      string
	Sum       string
}

func (b *BillJson) OverallPrice() float64 {
	return float64(b.Data.Json.TotalSum / 100)
}

func (b *BillJson) TransactionTime() (*time.Time, error) {
	timeString := b.Data.Json.DateTime
	time, err := time.Parse("2006-01-02T15:04:05", timeString)
	if err != nil {
		return nil, err
	}
	return &time, nil
}

func (b *BillJson) toBill() (*B.Bill, error) {
	// timeNow := time.Now()
	// timeTransaction, err := b.TransactionTime()
	// if err != nil {
	// return nil, err
	// }

	// bill := B.BillNew(
	// 	timeNow,
	// 	b.Data.Json.RetailPlace,
	// 	*timeTransaction,
	// 	b.OverallPrice(),
	// 	B.RUB,
	// )
	return &B.Bill{}, nil
}
