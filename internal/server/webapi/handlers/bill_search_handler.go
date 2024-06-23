package handlers

import (
	"billdb/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	BillsSearch = server.Get("/search/bills", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			return c.Render(http.StatusOK, "search-bills.html", map[string]interface{}{})
		}
	})

	BillSearchQueary = server.Post("/search/bills", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := make(map[string]interface{})
			r["success"] = false

			queryBase := `SELECT
				id, name, dates, price, currency, country, tag
			FROM
				bills
			WHERE
				name LIKE ?
				OR tag LIKE ?
				OR dates LIKE ?
				OR currency LIKE ?
				OR country LIKE ?
			ORDER BY
				dates DESC;`
			q := "%" + c.FormValue("q") + "%"

			c.Logger().Infof("query value: %s", q)

			db := s.BillRepo.GetDb()
			rows, err := db.Query(queryBase, q, q, q, q, q)
			if err != nil {
				r["message"] = "Error while querying the database"
				return c.Render(http.StatusOK, "search-bills-result.html", r)
			}
			defer rows.Close()

			var result []Bill
			for rows.Next() {
				var (
					Id       int64
					Name     string
					Date     string
					Price    float64
					Currency string
					Country  string
					Tag      string
				)
				rows.Scan(&Id, &Name, &Date, &Price, &Currency, &Country, &Tag)
				if err != nil {
					r["message"] = "Error while scanning the database"
					return c.Render(http.StatusOK, "search-bills-result.html", r)
				}
				result = append(result, Bill{
					Id:           Id,
					Name:         Name,
					Date:         Date,
					Price:        Price,
					Currency:     Currency,
					ExchangeRate: 0,
					Country:      Country,
					Tag:          Tag,
				})
			}

			r["result"] = result
			r["success"] = true
			return c.Render(http.StatusOK, "search-bills-result.html", r)
		}
	})
)
