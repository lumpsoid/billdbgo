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
				invoice.invoice_id, 
				item_name, 
				invoice_date, 
				item_price, 
				item_price_one, 
				item_quantity, 
				tag.tag_name
			FROM
				item
			LEFT JOIN invoice ON item.invoice_id = invoice.invoice_id
			LEFT JOIN item_tag ON item_tag.item_id = item.item_id
			LEFT JOIN tag ON tag.tag_id = item_tag.tag_id
			WHERE
				item_name LIKE ?
				OR tag.tag_name LIKE ?
				OR invoice_date LIKE ?
			ORDER BY
				invoice_date DESC;`
			q := "%" + c.FormValue("q") + "%"

			db := s.BillRepo.GetDb()
			rows, err := db.Query(queryBase, q, q, q)
			if err != nil {
				r["message"] = "Error while querying the database"
				return c.Render(http.StatusOK, "search-items-result.html", r)
			}
			defer rows.Close()

			var result []map[string]interface{}
			for rows.Next() {
				var (
					Id       string
					Name     string
					Date     string
					Price    float64
					PriceOne string
					Quantity string
					Tag      *string
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
