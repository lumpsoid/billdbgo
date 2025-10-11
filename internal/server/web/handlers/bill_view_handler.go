package handlers

import (
	"billdb/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	BillView = server.Get("/bill/:id", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Param("id")
			bill, err := s.BillRepo.GetBillByID(id)
			if err != nil {
				return err
			}
			items, err := s.BillRepo.GetItemsByID(id)
			if err != nil {
				return err
			}

			return c.Render(http.StatusOK, "bill-view.html", map[string]interface{}{
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
	})
)
