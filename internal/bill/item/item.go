package item

type Item struct {
	ItemId   string
	BillId   string
	Name     string
	Price    float64
	PriceOne float64
	Quantity float64
}

func New(
	itemId string,
	billId string,
	name string,
	price float64,
	priceOne float64,
	quantity float64,
) *Item {
	return &Item{
		ItemId:   itemId,
		BillId:   billId,
		Name:     name,
		Price:    price,
		PriceOne: priceOne,
		Quantity: quantity,
	}
}
