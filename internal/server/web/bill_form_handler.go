package web

import (
	"billdb/internal/bill"
	"billdb/internal/bill/country"
	"billdb/internal/bill/currency"
	"billdb/internal/bill/item"
	"billdb/internal/bill/tag"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
)

type BillRequest struct {
	Id           string  `form:"id"`
	Name         string  `form:"name"`
	Tag          string  `form:"tag"`
	Date         string  `form:"date"`
	Price        float64 `form:"price"`
	Currency     string  `form:"currency"`
	ExchangeRate float64 `form:"exchange_rate"`
	Country      string  `form:"country"`
}

func (w *WebHandlers) BillFormPage(c echo.Context) error {
	currencies := currency.Available()
	countries := country.Available()
	tags, err := w.BillRepo.GetTags()
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "bill-form.html", map[string]any{
		"currencies": currencies,
		"countries":  countries,
		"tags":       tags,
	})
}

func (w *WebHandlers) BillFormSubmit(c echo.Context) error {
	r := map[string]any{
		"success": false,
		"results": []map[string]any{},
	}
	b := new(BillRequest)
	responseHtml := "bill-insert-response.html"

	err := c.Bind(b)
	if err != nil {
		r["message"] = fmt.Sprintf("Error binding form data: %v", err)
		return c.Render(http.StatusOK, responseHtml, r)
	}
	b.Id = ksuid.New().String()

	// Create a single result entry for this form submission
	result := map[string]any{
		"link":    "Manual form entry",
		"success": false,
	}

	billCurrency, err := currency.Parse(b.Currency)
	if err != nil {
		result["message"] = fmt.Sprintf("Invalid currency: %v", err)
		r["results"] = append(r["results"].([]map[string]any), result)
		r["message"] = "Failed to process bill"
		return c.Render(http.StatusOK, responseHtml, r)
	}

	billCountry, err := country.Parse(b.Country)
	if err != nil {
		result["message"] = fmt.Sprintf("Invalid country: %v", err)
		r["results"] = append(r["results"].([]map[string]any), result)
		r["message"] = "Failed to process bill"
		return c.Render(http.StatusOK, responseHtml, r)
	}

	billDate, err := bill.StringToDate(b.Date)
	if err != nil {
		result["message"] = fmt.Sprintf("Invalid date format: %v", err)
		r["results"] = append(r["results"].([]map[string]any), result)
		r["message"] = "Failed to process bill"
		return c.Render(http.StatusOK, responseHtml, r)
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

	billDupCount, err := w.BillRepo.CheckDuplicateBill(billNew)
	if err != nil {
		result["message"] = fmt.Sprintf("Error checking duplicates: %v", err)
		r["results"] = append(r["results"].([]map[string]any), result)
		r["message"] = "Failed to process bill"
		return c.Render(http.StatusOK, responseHtml, r)
	}

	if billDupCount != 0 {
		result["message"] = "Found duplicate bills in database"
		result["dupInt"] = billDupCount
		r["results"] = append(r["results"].([]map[string]any), result)
		r["message"] = fmt.Sprintf("Found %d duplicate bill(s)", billDupCount)
		return c.Render(http.StatusOK, responseHtml, r)
	}

	err = w.BillRepo.InsertBill(billNew)
	if err != nil {
		result["message"] = fmt.Sprintf("Error inserting bill to database: %v", err)
		r["results"] = append(r["results"].([]map[string]any), result)
		r["message"] = "Failed to process bill"
		return c.Render(http.StatusOK, responseHtml, r)
	}

	billFromDb, err := w.BillRepo.GetBillByID(billNew.Id)
	if err != nil {
		result["message"] = fmt.Sprintf("Error retrieving bill from database: %v", err)
		r["results"] = append(r["results"].([]map[string]any), result)
		r["message"] = "Bill inserted but failed to retrieve"
		return c.Render(http.StatusOK, responseHtml, r)
	}

	result["success"] = true
	result["message"] = "Bill inserted successfully"
	result["bill"] = billFromDb
	r["results"] = append(r["results"].([]map[string]any), result)
	r["success"] = true
	r["message"] = "Bill processed successfully"

	return c.Render(http.StatusOK, responseHtml, r)
}
