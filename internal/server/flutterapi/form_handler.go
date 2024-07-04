package flutterapi

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
		r.Bill = make([]BillApi, 0)
		r.Force = false

		err := c.Bind(req)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusBadRequest, err)
		}
		r.Force = req.Force
		billDate, err := bill.StringToDate(req.Date)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusBadRequest, r)
		}
		billCountry, err := country.Parse(req.Country)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusBadRequest, r)
		}
		billCurrency, err := currency.Parse(req.Currency)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusBadRequest, r)
		}

		billId := ksuid.New()
		billAccepted := bill.New(
			billId.String(),
			req.Name,
			*billDate,
			req.Price,
			billCurrency,
			billCountry,
			[]*item.Item{},
			tag.New(req.Tags),
			"",
			"",
		)
		billApi := BillApi{
			Id:           billAccepted.Id,
			Name:         req.Name,
			Date:         billAccepted.GetDateString(),
			Price:        req.Price,
			Currency:     req.Currency,
			ExchangeRate: req.ExchangeRate,
			Country:      req.Country,
			Items:        0,
			Link:         "",
			Duplicates:   0,
		}

		billDupCount, err := s.BillRepo.CheckDuplicateBill(billAccepted)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusInternalServerError, r)
		}
		if billDupCount != 0 {
			r.Success = "duplicates"
			billApi.Duplicates = billDupCount
			r.Message = fmt.Sprintf("Find duplicates in the db = %d\n", billDupCount)
			return c.JSON(http.StatusOK, r)
		}

		err = s.BillRepo.InsertBill(billAccepted)
		if err != nil {
			r.Message = fmt.Sprintf("%v", err)
			return c.JSON(http.StatusInternalServerError, r)
		}

		r.Success = "success"
		r.Bill = []BillApi{billApi}

		return c.JSON(http.StatusOK, r)
	}
})
