package api

import (
	"billdb/internal/parser"
	"billdb/internal/server"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestQr struct {
	Link  string `json:"link"`
	Force bool   `json:"force"`
}

var QrHandler = server.Post(baseApiPath+"/qr", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(RequestQr)
		r := new(ResponseFlutter)
		r.Success = "error"
		r.Bill = make([]BillApi, 0)

		if err := c.Bind(req); err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusBadRequest, r)
		}
		r.Force = req.Force

		if req.Link == "" {
			r.Message = "Empty link"
			return c.JSON(http.StatusBadRequest, r)
		}

		p, err := parser.GetBillParser(req.Link)
		if err != nil {
			r.Message = fmt.Sprintf("Error while getting parser for the url: %v", err)
			return c.JSON(http.StatusInternalServerError, r)
		}

		bill, err := p.Parse(req.Link)
		if err != nil {
			r.Message = "Error while parsing the site"
			return c.JSON(http.StatusInternalServerError, r)
		}
		b := BillApi{
			// TODO check with app, what If I will send string in timestamp
			Id:       bill.Id,
			Name:     bill.Name,
			Date:     bill.GetDateString(),
			Price:    bill.Price,
			Currency: bill.GetCurrencyString(),
			// TODO exchange rate system
			// ExchangeRate: bill.ExchangeRate,
			Country: bill.GetCountryString(),
			Items:   len(bill.Items),
			Link:    req.Link,
		}
		r.Bill = []BillApi{b}

		// TODO check in flutter app, do I need to send beck duplicates?
		billDupCount, err := s.BillRepo.CheckDuplicateBill(bill)
		if err != nil {
			r.Message = fmt.Sprintf("Duplicates error: %v", err)
			return c.JSON(http.StatusInternalServerError, r)
		}
		if billDupCount != 0 {
			r.Success = "duplicates"
			b.Duplicates = billDupCount
			r.Message = fmt.Sprintf("Find duplicates in the db = %d\n", billDupCount)
			return c.JSON(http.StatusOK, r)
		}

		err = s.BillRepo.InsertBill(bill)
		if err != nil {
			r.Message = fmt.Sprintf("Error while inserting a bill: %v", err)
			return c.JSON(http.StatusInternalServerError, r)
		}

		err = s.BillRepo.InsertItems(bill.Items)
		if err != nil {
			r.Message = fmt.Sprintf("Error while inserting items: %v", err)
			return c.JSON(http.StatusInternalServerError, r)
		}

		r.Success = "success"
		r.Bill = []BillApi{b}

		return c.JSON(http.StatusOK, r)
	}
})
