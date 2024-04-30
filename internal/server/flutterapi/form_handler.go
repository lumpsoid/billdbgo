package flutterapi

import (
	B "billdb/internal/bill"
	"billdb/internal/server"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type RequestForm struct {
	Name         string  `json:"name"`
	Date         string  `json:"date"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	ExchangeRate float64 `json:"exchange_rate"`
	Country      string  `json:"country"`
	Tags         string  `json:"tags"`
	Force        bool    `json:"force"`
}

var FormHandler = server.Post(baseApiPath+"/form", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(RequestForm)
		r := new(ResponseFlutter)
		r.Success = "error"
		r.Bill = make([]Bill, 0)
		r.Force = false

		err := c.Bind(req)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusBadRequest, err)
		}
		r.Force = req.Force
		billDate, err := B.StringToDate(req.Date)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusBadRequest, r)
		}
		billCountry, err := B.StringToCountry(req.Country)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusBadRequest, r)
		}
		billCurrency, err := B.StringToCurrency(req.Currency)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusBadRequest, r)
		}

		bill := B.BillNew(
			time.Now().Local(),
			req.Name,
			*billDate,
			req.Price,
			billCurrency,
			req.ExchangeRate,
			billCountry,
			[]*B.Item{},
			B.Tag(req.Tags),
			"",
			"",
		)
		b := Bill{
			Timestamp:    bill.GetIdUnix(),
			Name:         req.Name,
			Date:         bill.GetDateString(),
			Price:        req.Price,
			Currency:     req.Currency,
			ExchangeRate: req.ExchangeRate,
			Country:      req.Country,
			Items:        0,
			Link:         "",
			Duplicates:   0,
		}

		billsDup, err := s.BillRepo.CheckDuplicateBill(bill)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusInternalServerError, r)
		}
		if len(billsDup) != 0 {
			r.Success = "duplicates"
			b.Duplicates = len(billsDup)
			r.Message = fmt.Sprintf("Find duplicates in the db = %d\n", len(billsDup))
			return c.JSON(http.StatusOK, r)
		}

		err = s.BillRepo.InsertBill(bill)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusInternalServerError, r)
		}

		r.Success = "success"
		r.Bill = []Bill{b}

		return c.JSON(http.StatusOK, r)
	}
})
