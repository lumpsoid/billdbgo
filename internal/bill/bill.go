package bill

import (
	"strconv"
	"strings"
	"time"
)

type Bill struct {
	Id           time.Time
	Name         string
	Date         time.Time
	Price        float64
	Currency     Currency
	ExchangeRate float64
	Country      Country
	Items        []*Item
	Tag          Tag
	Link         string
	BillText     string
}

func BillNew(
	id time.Time,
	name string,
	date time.Time,
	price float64,
	currency Currency,
	exchangeRate float64,
	country Country,
	items []*Item,
	tag Tag,
	link string,
	billText string,
) *Bill {
	return &Bill{
		Id:           id,
		Name:         name,
		Date:         date,
		Price:        price,
		Currency:     currency,
		ExchangeRate: exchangeRate,
		Country:      country,
		Items:        items,
		Tag:          tag,
		Link:         link,
		BillText:     billText,
	}
}

func (b *Bill) AddItem(item *Item) {
	b.Items = append(b.Items, item)
}

func (b *Bill) GetIdUnix() int64 {
	return b.Id.Local().UnixNano()
}

func (b *Bill) GetDateString() string {
	return b.Date.Format("2006-01-02")
}

func (b *Bill) GetCurrencyString() string {
	return currencyToString[b.Currency]
}

func (b *Bill) GetCountryString() string {
	return b.Country.String()
}

func UpdateBillProperty(bill *Bill, property string, value interface{}) error {
	switch property {
	case "name":
		bill.Name = strings.Trim(value.(string), " ,\t-")
	case "date":
		dateNew, err := StringToDate(value.(string))
		if err != nil {
			return err
		}
		bill.Date = *dateNew
	case "price":
		priceNew, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return err
		}
		bill.Price = priceNew
	case "currency":
		currencyNew, err := StringToCurrency(value.(string))
		if err != nil {
			return err
		}
		bill.Currency = currencyNew
	case "exchange_rate":
		exchangeRateNew, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			return err
		}
		bill.ExchangeRate = exchangeRateNew
	case "country":
		countryNew, err := StringToCountry(value.(string))
		if err != nil {
			return err
		}
		bill.Country = countryNew
	case "tag":
		bill.Tag = Tag(value.(string))
	case "link":
		bill.Link = value.(string)
	}
	return nil
}

func DateToString(date time.Time) string {
	return date.Format("2006-01-02")
}

func StringToDate(date string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func IdToUnix(id time.Time) int64 {
	return id.Local().UnixMilli()
}

func UnixToId(id int64) time.Time {
	return time.Unix(0, id)
}
