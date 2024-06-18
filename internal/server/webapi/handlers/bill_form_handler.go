package handlers

import (
	B "billdb/internal/bill"
	"billdb/internal/server"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Bill struct {
	Id           int64   `form:"id"`
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
			return c.Render(http.StatusOK, "bill-form.html", nil)
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
			b.Id = time.Now().Local().UnixNano()

			billCurrency, err := B.StringToCurrency(b.Currency)
			if err != nil {
				r["success"] = false
				r["message"] = fmt.Sprintf("%v", err)
				return c.Render(http.StatusOK, "bill-form-response.html", r)
			}
			billCountry, err := B.StringToCountry(b.Country)
			if err != nil {
				r["success"] = false
				r["message"] = fmt.Sprintf("%v", err)
				return c.Render(http.StatusOK, "bill-form-response.html", r)
			}
			billDate, err := B.StringToDate(b.Date)
			if err != nil {
				r["success"] = false
				r["message"] = fmt.Sprintf("%v", err)
				return c.Render(http.StatusOK, "bill-form-response.html", r)
			}

			bill := B.BillNew(
				B.UnixToId(b.Id),
				b.Name,
				*billDate,
				b.Price,
				billCurrency,
				b.ExchangeRate,
				billCountry,
				[]*B.Item{},
				B.Tag(b.Tag),
				"",
				"",
			)

			billsDup, err := s.BillRepo.CheckDuplicateBill(bill)
			if err != nil {
				r["success"] = false
				r["message"] = fmt.Sprintf("%v", err)
				return c.Render(http.StatusOK, "bill-form-response.html", r)
			}
			if len(billsDup) != 0 {
				r["success"] = false
				r["message"] = fmt.Sprintf("Find duplicates in the db = %d", len(billsDup))
				r["dupId"] = billsDup[0].GetIdUnix()
				return c.Render(http.StatusOK, "bill-form-response.html", r)
			}

			err = s.BillRepo.InsertBill(bill)
			if err != nil {
				r["success"] = false
				r["message"] = "Error while inserting bill to db"
				return c.Render(http.StatusOK, "bill-form-response.html", r)
			}
			r["success"] = true
			r["message"] = "Bill parsed successfully"
			r["bill"] = b
			return c.Render(http.StatusOK, "bill-form-response.html", r)
		}
	})
)
