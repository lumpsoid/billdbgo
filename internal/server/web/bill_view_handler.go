package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (w *WebHandlers) BillView(c echo.Context) error {
	id := c.Param("id")
	bill, err := w.BillRepo.GetBillByID(id)
	if err != nil {
		return err
	}
	items, err := w.BillRepo.GetItemsByID(id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "bill-view.html", map[string]any{
		"id":       bill.Id,
		"date":     bill.GetDateString(),
		"name":     bill.Name,
		"price":    bill.Price,
		"currency": bill.GetCurrencyString(),
		// TODO exchange rate system
		// "exchange_rate": bill.ExchangeRate,
		"country":   bill.GetCountryString(),
		"tag":       bill.Tag.String,
		"link":      bill.Link,
		"bill_text": bill.BillText,
		"items":     items,
	})
}
