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
	BillBrowseLanding = server.Get("/bills", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			timeNow := time.Now()
			return c.Redirect(http.StatusMovedPermanently, "/bills/"+timeNow.Format("2006/01"))
		}
	})

	BillBrowse = server.Get("/bills/:y/:m", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := make(map[string]interface{})
			r["success"] = false

			year, err := strconv.ParseInt(c.Param("y"), 10, 64)
			if err != nil {
				r["message"] = fmt.Sprintf("Invalid year: %s | URL: %s", c.Param("y"), c.Request().URL)
				return c.Render(http.StatusOK, "bills-browse.html", r)
			}
			month, err := strconv.ParseInt(c.Param("m"), 10, 64)
			if err != nil {
				r["message"] = fmt.Sprintf("Invalid month: %s | URL: %s", c.Param("m"), c.Request().URL)
				return c.Render(http.StatusOK, "bills-browse.html", r)
			}

			timeNow := time.Now()
			timeRequested := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			if timeRequested.After(timeNow) {
				r["message"] = "Requested date is in the future"
				return c.Render(http.StatusOK, "bills-browse.html", r)
			}

			queryBase := `
      SELECT
				id, name, dates, price, currency, exchange_rate, country, tag
			FROM
				bills
			WHERE
				strftime('%%Y-%%m', dates) = '%d-%02d'
			ORDER BY
				dates DESC;
      `
			query := fmt.Sprintf(queryBase, year, month)

			db := s.BillRepo.GetDb()
			rows, err := db.Query(query)
			if err != nil {
				r["message"] = "Error while querying the database"
				return c.Render(http.StatusOK, "bills-browse.html", r)
			}
			defer rows.Close()

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
					r["message"] = "Error while scanning the database"
					return c.Render(http.StatusOK, "bills-browse.html", r)
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

			nextMonth := timeRequested.AddDate(0, 1, 0)
			if nextMonth.Before(timeNow) {
				r["nextPage"] = c.Echo().Reverse("bills", nextMonth.Year(), int(nextMonth.Month()))
			}
			prevMonth := timeRequested.AddDate(0, -1, 0)
			r["prevPage"] = c.Echo().Reverse("bills", prevMonth.Year(), int(prevMonth.Month()))

			r["success"] = true
			r["bills"] = billsResponse
			return c.Render(http.StatusOK, "bills-browse.html", r)
		}
	})
)
