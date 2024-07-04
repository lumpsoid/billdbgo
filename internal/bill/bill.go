package bill

import (
	"billdb/internal/bill/country"
	"billdb/internal/bill/currency"
	"billdb/internal/bill/item"
	"billdb/internal/bill/tag"
	"strconv"
	"strings"
	"time"
)

type Bill struct {
	Id       string
	Name     string
	Date     time.Time
	Price    float64
	Currency currency.Currency
	Country  country.Country
	Items    []*item.Item
	Tag      *tag.Tag // TODO add posibility to be nil
	Link     string
	BillText string // TODO transform into a struct
}

func New(
	id string,
	name string,
	date time.Time,
	price float64,
	currency currency.Currency,
	country country.Country,
	items []*item.Item,
	tag *tag.Tag,
	link string,
	billText string,
) *Bill {
	return &Bill{
		Id:       id,
		Name:     name,
		Date:     date,
		Price:    price,
		Currency: currency,
		Country:  country,
		Items:    items,
		Tag:      tag,
		Link:     link,
		BillText: billText,
	}
}

func (b *Bill) AddItem(item *item.Item) {
	b.Items = append(b.Items, item)
}

func (b *Bill) GetDateString() string {
	return b.Date.Format("2006-01-02")
}

func (b *Bill) GetCurrencyString() string {
	return b.Currency.String()
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
		currencyNew, err := currency.Parse(value.(string))
		if err != nil {
			return err
		}
		bill.Currency = currencyNew
	case "exchange_rate":
		// TODO migrate exchange rate system
		// exchangeRateNew, err := strconv.ParseFloat(value.(string), 64)
		// if err != nil {
		// return err
		// }
		// bill.ExchangeRate = exchangeRateNew
	case "country":
		countryNew, err := country.Parse(value.(string))
		if err != nil {
			return err
		}
		bill.Country = countryNew
	case "tag":
		bill.Tag = tag.New(value.(string))
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
