package handlers

import (
	"billdb/internal/parser"
	"billdb/internal/server"
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool
	Message string
	Bill    Bill
}

var BillFromLink = server.Get("/bill-from-link", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.Render(http.StatusOK, "bill-from-link.html", nil)
	}
})

var BillFromLinkResponse = server.Post("/bill-from-link", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		var r Response
		link := c.FormValue("link")

		u, err := url.Parse(link)
		if err != nil {
			c.Logger().Errorf("Error while parsing link from user: %v\n", err)
			r.Success = false
			r.Message = "Not valid URL"
			return c.Render(http.StatusOK, "bill-from-link-response.html", r)
		}
		p, err := parser.GetBillParser(u)
		if err != nil {
			r.Success = false
			r.Message = fmt.Sprintf("%v\n", err)
			return c.Render(http.StatusOK, "bill-from-link-response.html", r)
		}
		bill, err := p.Parse(u)
		if err != nil {
			r.Success = false
			r.Message = fmt.Sprintf("%v\n", err)
			return c.Render(http.StatusOK, "bill-from-link-response.html", r)
		}

		billDupCount, err := s.BillRepo.CheckDuplicateBill(bill)
		if err != nil {
			r.Success = false
			r.Message = fmt.Sprintf("%v\n", err)
			return c.Render(http.StatusOK, "bill-from-link-response.html", r)
		}
		if billDupCount != -1 {
			r.Success = false
			r.Message = fmt.Sprintf("Find duplicates in the db = %d\n", billDupCount)
			return c.Render(http.StatusOK, "bill-from-link-response.html", r)
		}

		err = s.BillRepo.InsertBill(bill)
		if err != nil {
			r.Success = false
			r.Message = "Error while inserting bill to db"
			return c.Render(http.StatusOK, "bill-from-link-response.html", r)
		}
		err = s.BillRepo.InsertItems(bill.Items)
		if err != nil {
			r.Success = false
			r.Message = "Error while inserting items to db"
			return c.Render(http.StatusOK, "bill-from-link-response.html", r)
		}
		r.Success = true
		r.Message = "Bill parsed successfully"
		r.Bill = Bill{
			Id:       bill.Id,
			Name:     bill.Name,
			Tag:      bill.Tag.String(),
			Date:     bill.GetDateString(),
			Price:    bill.Price,
			Currency: bill.GetCurrencyString(),
			// TODO exchange rate system
			// ExchangeRate: bill.ExchangeRate,
			Country: bill.GetCountryString(),
		}

		return c.Render(http.StatusOK, "bill-from-link-response.html", r)
	}
},
)
