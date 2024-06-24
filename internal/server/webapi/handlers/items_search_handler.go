package handlers

import (
	"billdb/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ItemsSearch = server.Get("/search/items", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			return c.Render(http.StatusOK, "search-items.html", map[string]interface{}{})
		}
	})

	ItemsSearchQueary = server.Post("/search/items", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := make(map[string]interface{})
			r["success"] = false

			queryBase := `SELECT
				i.id, i.name, b.dates, i.price, i.price_one, i.quantity, m.tag
			FROM
				items as i
			LEFT JOIN bills as b ON i.id = b.id
			LEFT JOIN items_meta as m ON m.name = i.name
			WHERE
			i.name LIKE ?
			OR m.tag LIKE ?
			OR b.dates LIKE ?
			ORDER BY
				dates DESC;`
			q := "%" + c.FormValue("q") + "%"

			db := s.BillRepo.GetDb()
			rows, err := db.Query(queryBase, q, q, q, q, q)
			if err != nil {
				r["message"] = "Error while querying the database"
				return c.Render(http.StatusOK, "search-items-result.html", r)
			}
			defer rows.Close()

			var result []map[string]interface{}
			for rows.Next() {
				var (
					Id       int64
					Name     string
					Date     string
					Price    float64
					PriceOne string
					Quantity string
					Tag      string
				)
				rows.Scan(&Id, &Name, &Date, &Price, &PriceOne, &Quantity, &Tag)
				if err != nil {
					r["message"] = "Error while scanning the database"
					return c.Render(http.StatusOK, "search-items-result.html", r)
				}
				result = append(result, map[string]interface{}{
					"Id":       Id,
					"Name":     Name,
					"Date":     Date,
					"Price":    Price,
					"PriceOne": PriceOne,
					"Quantity": Quantity,
					"Tag":      Tag,
				})
			}

			r["result"] = result
			r["success"] = true
			return c.Render(http.StatusOK, "search-items-result.html", r)
		}
	})
)
