package bill

import "time"

type Item struct {
	Id       time.Time
	Name     string
	Price    float64
	PriceOne float64
	Quantity float64
}

func ItemNew(
	id time.Time,
	name string,
	price float64,
	priceOne float64,
	quantity float64,
) *Item {
	return &Item{
		Id:       id,
		Name:     name,
		Price:    price,
		PriceOne: priceOne,
		Quantity: quantity,
	}
}

func (t *Item) GetIdUnix() int64 {
  return t.Id.Local().UnixNano()
}
