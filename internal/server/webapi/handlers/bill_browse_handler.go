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
	BillBrowseLanding = server.Get("/browse/bills", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			timeNow := time.Now()
			return c.Redirect(http.StatusMovedPermanently, "/browse/bills/"+timeNow.Format("2006/01"))
		}
	})

	BillBrowse = server.Get("/browse/bills/:y/:m", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := make(map[string]interface{})
			r["success"] = false

			year, err := strconv.ParseInt(c.Param("y"), 10, 64)
			if err != nil {
				r["message"] = fmt.Sprintf("Invalid year: %s | URL: %s", c.Param("y"), c.Request().URL)
				return c.Render(http.StatusOK, "browse-bills.html", r)
			}
			month, err := strconv.ParseInt(c.Param("m"), 10, 64)
			if err != nil {
				r["message"] = fmt.Sprintf("Invalid month: %s | URL: %s", c.Param("m"), c.Request().URL)
				return c.Render(http.StatusOK, "browse-bills.html", r)
			}

			timeNow := time.Now()
			timeRequested := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			if timeRequested.After(timeNow) {
				r["message"] = "Requested date is in the future"
				return c.Render(http.StatusOK, "browse-bills.html", r)
			}

			queryBase := `SELECT
				invoice.invoice_id, 
				invoice_name, 
				invoice_date, 
				invoice_price, 
				invoice_currency, 
				invoice_country,
				tag.tag_name
			FROM
				invoice
			LEFT JOIN invoice_tag ON invoice_tag.invoice_id = invoice.invoice_id
			LEFT JOIN tag ON tag.tag_id = invoice_tag.tag_id
			WHERE
				strftime('%%Y-%%m', invoice_date) = '%d-%02d'
			ORDER BY
				invoice_date DESC;`
			query := fmt.Sprintf(queryBase, year, month)

			db := s.BillRepo.GetDb()
			rows, err := db.Query(query)
			if err != nil {
				r["message"] = fmt.Sprintf("Error while querying the database: %v; Db path: %s", err, s.Config.DbPath)
				return c.Render(http.StatusOK, "browse-bills.html", r)
			}
			defer rows.Close()

			var billsResponse []Bill
			for rows.Next() {
				var (
					Id       string
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
					return c.Render(http.StatusOK, "browse-bills.html", r)
				}
				billsResponse = append(billsResponse, Bill{
					Id:       Id,
					Name:     Name,
					Date:     Date,
					Price:    Price,
					Currency: Currency,
					// TODO echange rate
					ExchangeRate: 0,
					Country:      Country,
					Tag:          Tag,
				})
			}

			nextMonth := timeRequested.AddDate(0, 1, 0)
			if nextMonth.Before(timeNow) {
				r["nextPage"] = c.Echo().Reverse("browse-bills", nextMonth.Year(), int(nextMonth.Month()))
			}
			prevMonth := timeRequested.AddDate(0, -1, 0)
			r["prevPage"] = c.Echo().Reverse("browse-bills", prevMonth.Year(), int(prevMonth.Month()))

			r["bills"] = billsResponse
			r["year"] = year
			r["month"] = fmt.Sprintf("%02d", month)
			r["success"] = true
			return c.Render(http.StatusOK, "browse-bills.html", r)
		}
	})
)
