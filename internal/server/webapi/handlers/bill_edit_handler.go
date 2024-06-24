package handlers

import (
	B "billdb/internal/bill"
	"billdb/internal/server"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var (
	BillEditPage = server.Get("/bill/:id/edit", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			billId, err := strconv.ParseInt(c.Param("id"), 10, 64)
			if err != nil {
				return err
			}
			bill, err := s.BillRepo.GetBillByID(billId)
			if err != nil {
				return err
			}
			currencies := B.GetCurrencyList()
			countries := B.GetCountryList()
			tags, err := s.BillRepo.GetTags()
			if err != nil {
				return err
			}

			return c.Render(http.StatusOK, "bill-edit.html", map[string]interface{}{
				"id":            bill.GetIdUnix(),
				"date":          bill.GetDateString(),
				"name":          bill.Name,
				"price":         bill.Price,
				"currency":      bill.GetCurrencyString(),
				"exchange_rate": bill.ExchangeRate,
				"country":       bill.GetCountryString(),
				"tag":           bill.Tag.String(),
				"link":          bill.Link,
				"bill_text":     bill.BillText,
				"currencies":    currencies,
				"countries":     countries,
				"tags":          tags,
			})
		}
	})

	BillEditSubmit = server.Put("/bill/:id/edit", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := make(map[string]interface{})
			r["success"] = false
			billId, err := strconv.ParseInt(c.Param("id"), 10, 64)
			if err != nil {
				c.Logger().Errorf("Error parsing id from url path: %v", err)
				r["error"] = "Error parsing id from url path"
				return c.Render(
					http.StatusOK,
					"bill-edit-result.html",
					r,
				)
			}
			bill, err := s.BillRepo.GetBillByID(billId)
			if err != nil {
				c.Logger().Errorf("Error getting bill by id: %v", err)
				return err
			}
			r["id"] = bill.GetIdUnix()
			r["cDate"] = bill.GetDateString()
			r["cName"] = bill.Name
			r["cPrice"] = bill.Price
			r["cCurrency"] = bill.GetCurrencyString()
			r["cExchangeRate"] = bill.ExchangeRate
			r["cCountry"] = bill.GetCountryString()
			r["cTag"] = bill.Tag.String()
			r["cLink"] = bill.Link

			params, err := c.FormParams()
			if err != nil {
				c.Logger().Errorf("Error reading form params: %v", err)
				return err
			}
			for property, value := range params {
				if len(value) == 0 {
					continue
				}
				if value[0] == "" {
					continue
				}
				err := B.UpdateBillProperty(bill, property, value[0])
				if err != nil {
					c.Logger().Errorf("Error wile updating bill property: %v", err)
					r["error"] = err
					return c.Render(
						http.StatusOK,
						"bill-edit-result.html",
						r,
					)
				}
			}
			err = s.BillRepo.UpdateBill(bill)
			if err != nil {
				c.Logger().Errorf("Error updating bill: %v", err)
				r["error"] = "Error updating bill in db."
				return c.Render(
					http.StatusOK,
					"bill-edit-result.html",
					r,
				)
			}
			r["nDate"] = bill.GetDateString()
			r["nName"] = bill.Name
			r["nPrice"] = bill.Price
			r["nCurrency"] = bill.GetCurrencyString()
			r["nExchangeRate"] = bill.ExchangeRate
			r["nCountry"] = bill.GetCountryString()
			r["nTag"] = bill.Tag.String()
			r["nLink"] = bill.Link

			c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			r["success"] = true
			return c.Render(
				http.StatusOK,
				"bill-edit-result.html",
				r,
			)
		}
	})
)
