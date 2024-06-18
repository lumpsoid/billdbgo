package handlers

import (
	"billdb/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	BillBrowse = server.Get("/bills", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {

			db := s.BillRepo.GetDb()

			query := `
      SELECT
        id, name, dates, price, currency, exchange_rate, country, tag
      FROM bills
      ORDER BY dates DESC
      LIMIT 20;
      `
			rows, err := db.Query(query)
			if err != nil {
				return err
			}

			var billsResponse []Bill
			for rows.Next() {
				var (
					Id           int64
					Name         string
					Date         string
					Price        float64
					Currency     string
					ExchangeRate float64
					Country      string
					Tag          string
				)
				rows.Scan(&Id, &Name, &Date, &Price, &Currency, &ExchangeRate, &Country, &Tag)
				if err != nil {
					return err
				}
				billsResponse = append(billsResponse, Bill{
					Id:           Id,
					Name:         Name,
					Date:         Date,
					Price:        Price,
					Currency:     Currency,
					ExchangeRate: ExchangeRate,
					Country:      Country,
					Tag:          Tag,
				})
			}

			return c.Render(http.StatusOK, "bills-browse.html", map[string]interface{}{
				"bills": billsResponse,
			})
		}
	})
)
