package handlers

import (
	"billdb/internal/server"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	ItemsBrowse = server.Get("/browse/items/:y/:m", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := make(map[string]interface{})
			r["success"] = false

			year, err := strconv.ParseInt(c.Param("y"), 10, 64)
			if err != nil {
				r["message"] = fmt.Sprintf("Invalid year: %s | URL: %s", c.Param("y"), c.Request().URL)
				return c.Render(http.StatusOK, "browse-items.html", r)
			}
			month, err := strconv.ParseInt(c.Param("m"), 10, 64)
			if err != nil {
				r["message"] = fmt.Sprintf("Invalid month: %s | URL: %s", c.Param("m"), c.Request().URL)
				return c.Render(http.StatusOK, "browse-items.html", r)
			}

			timeNow := time.Now()
			timeRequested := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			if timeRequested.After(timeNow) {
				r["message"] = "Requested date is in the future"
				return c.Render(http.StatusOK, "browse-items.html", r)
			}

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
				strftime('%%Y-%%m', invoice_date) = '%d-%02d'
			ORDER BY
				invoice_date DESC;`
			query := fmt.Sprintf(queryBase, year, month)

			db := s.BillRepo.GetDb()
			rows, err := db.Query(query)
			if err != nil {
				r["message"] = "Error while querying the database"
				return c.Render(http.StatusOK, "browse-items.html", r)
			}
			defer rows.Close()

			var itemsResponse []map[string]interface{}
			for rows.Next() {
				var (
					Id       string
					Name     string
					Date     string
					Price    float64
					PriceOne string
					Quantity float64
					Tag      *string
				)
				rows.Scan(&Id, &Name, &Date, &Price, &PriceOne, &Quantity, &Tag)
				if err != nil {
					r["message"] = "Error while scanning the database"
					return c.Render(http.StatusOK, "browse-items.html", r)
				}
				itemsResponse = append(itemsResponse, map[string]interface{}{
					"Id":       Id,
					"Name":     Name,
					"Date":     Date,
					"Price":    Price,
					"PriceOne": PriceOne,
					"Quantity": Quantity,
					"Tag":      Tag,
				})
			}

			nextMonth := timeRequested.AddDate(0, 1, 0)
			if nextMonth.Before(timeNow) {
				r["nextPage"] = c.Echo().Reverse("browse-items", nextMonth.Year(), int(nextMonth.Month()))
			}
			prevMonth := timeRequested.AddDate(0, -1, 0)
			r["prevPage"] = c.Echo().Reverse("browse-items", prevMonth.Year(), int(prevMonth.Month()))

			r["items"] = itemsResponse
			r["CurrentMonthItemsPage"] = "/browse/bills/" + timeNow.Format("2006/01")
			r["year"] = year
			r["month"] = fmt.Sprintf("%02d", month)
			r["success"] = true
			return c.Render(http.StatusOK, "browse-items.html", r)
		}
	})
)
