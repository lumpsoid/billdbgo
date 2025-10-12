package web

import (
	"billdb/internal/bill"
	"billdb/internal/bill/country"
	"billdb/internal/bill/currency"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (w *WebHandlers) BillEditPage(c echo.Context) error {
	billId := c.Param("id")
	billRequested, err := w.BillRepo.GetBillByID(billId)
	if err != nil {
		return err
	}
	currencies := currency.Available()
	countries := country.Available()
	tags, err := w.BillRepo.GetTags()
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "bill-edit.html", map[string]any{
		"id":       billRequested.Id,
		"date":     billRequested.GetDateString(),
		"name":     billRequested.Name,
		"price":    billRequested.Price,
		"currency": billRequested.GetCurrencyString(),
		// TODO exchange rate system
		// "exchange_rate": bill.ExchangeRate,
		"country":    billRequested.GetCountryString(),
		"tag":        billRequested.Tag.String,
		"link":       billRequested.Link,
		"bill_text":  billRequested.BillText,
		"currencies": currencies,
		"countries":  countries,
		"tags":       tags,
	})
}

func (w *WebHandlers) BillEditSubmit(c echo.Context) error {
	r := make(map[string]interface{})
	r["success"] = false
	billId := c.Param("id")
	billEdited, err := w.BillRepo.GetBillByID(billId)
	if err != nil {
		c.Logger().Errorf("Error getting bill by id: %v", err)
		return err
	}
	r["id"] = billEdited.Id
	r["cDate"] = billEdited.GetDateString()
	r["cName"] = billEdited.Name
	r["cPrice"] = billEdited.Price
	r["cCurrency"] = billEdited.GetCurrencyString()
	// TODO exchange rate system
	// r["cExchangeRate"] = bill.ExchangeRate
	r["cCountry"] = billEdited.GetCountryString()
	r["cTag"] = billEdited.Tag.String
	r["cLink"] = billEdited.Link

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
		err := bill.UpdateBillProperty(billEdited, property, value[0])
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
	err = w.BillRepo.UpdateBill(billEdited)
	if err != nil {
		c.Logger().Errorf("Error updating bill: %v", err)
		r["error"] = "Error updating bill in db."
		return c.Render(
			http.StatusOK,
			"bill-edit-result.html",
			r,
		)
	}
	billNew, err := w.BillRepo.GetBillByID(billId)
	if err != nil {
		c.Logger().Errorf("Error getting bill by id: %v", err)
		return err
	}
	r["nDate"] = billNew.GetDateString()
	r["nName"] = billNew.Name
	r["nPrice"] = billNew.Price
	r["nCurrency"] = billNew.GetCurrencyString()
	// TODO exchange rate system
	// r["nExchangeRate"] = bill.ExchangeRate
	r["nCountry"] = billNew.GetCountryString()
	r["nTag"] = billNew.Tag.String
	r["nLink"] = billNew.Link

	// c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	r["success"] = true
	return c.Render(
		http.StatusOK,
		"bill-edit-result.html",
		r,
	)
}
