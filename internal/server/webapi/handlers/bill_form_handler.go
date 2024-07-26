package handlers

import (
	"billdb/internal/bill"
	"billdb/internal/bill/country"
	"billdb/internal/bill/currency"
	"billdb/internal/bill/item"
	"billdb/internal/bill/tag"
	"billdb/internal/server"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
)

type Bill struct {
	Id           string  `form:"id"`
	Name         string  `form:"name"`
	Tag          string  `form:"tag"`
	Date         string  `form:"date"`
	Price        float64 `form:"price"`
	Currency     string  `form:"currency"`
	ExchangeRate float64 `form:"exchange_rate"`
	Country      string  `form:"country"`
}

var (
	BillFormPage = server.Get("/bill/form", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			currencies := currency.Available()
			countries := country.Available()
			tags, err := s.BillRepo.GetTags()
			if err != nil {
				return err
			}
			return c.Render(http.StatusOK, "bill-form.html", map[string]interface{}{
				"currencies": currencies,
				"countries":  countries,
				"tags":       tags,
			})
		}
	})

	BillFormSubmit = server.Post("/bill/form", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := make(map[string]interface{})
			b := new(Bill)

			err := c.Bind(b)
			if err != nil {
				return err
			}
			b.Id = ksuid.New().String()

			billCurrency, err := currency.Parse(b.Currency)
			if err != nil {
				r["success"] = false
				r["message"] = fmt.Sprintf("%v", err)
				return c.Render(http.StatusOK, "bill-insert-response.html", r)
			}
			billCountry, err := country.Parse(b.Country)
			if err != nil {
				r["success"] = false
				r["message"] = fmt.Sprintf("%v", err)
				return c.Render(http.StatusOK, "bill-insert-response.html", r)
			}
			billDate, err := bill.StringToDate(b.Date)
			if err != nil {
				r["success"] = false
				r["message"] = fmt.Sprintf("%v", err)
				return c.Render(http.StatusOK, "bill-insert-response.html", r)
			}

			billNew := bill.New(
				b.Id,
				b.Name,
				*billDate,
				b.Price,
				billCurrency,
				billCountry,
				[]*item.Item{},
				tag.New(b.Tag),
				"",
				"",
			)

			billDupCount, err := s.BillRepo.CheckDuplicateBill(billNew)
			if err != nil {
				r["success"] = false
				r["message"] = fmt.Sprintf("%v", err)
				return c.Render(http.StatusOK, "bill-insert-response.html", r)
			}
			if billDupCount != 0 {
				r["success"] = false
				r["message"] = fmt.Sprintf("Find duplicates in the db = %d", billDupCount)
        // TODO dupId implement to return ids of the duplicate bills 
				r["dupId"] = "test"
				return c.Render(http.StatusOK, "bill-insert-response.html", r)
			}

			err = s.BillRepo.InsertBill(billNew)
			if err != nil {
				r["success"] = false
				r["message"] = "Error while inserting bill to db"
				return c.Render(http.StatusOK, "bill-insert-response.html", r)
			}
			r["success"] = true
			r["message"] = "Bill parsed successfully"
			r["bill"] = b
			return c.Render(http.StatusOK, "bill-insert-response.html", r)
		}
	})
)
