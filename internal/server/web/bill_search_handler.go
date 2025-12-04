package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (w *WebHandlers) BillsSearch(c echo.Context) error {
	return c.Render(http.StatusOK, "search-bills.html", nil)
}

func (w *WebHandlers) BillSearchQueary(c echo.Context) error {
	r := make(map[string]any)
	r["success"] = false

	queryBase := `SELECT
				invoice.invoice_id, 
				invoice_name, 
				invoice_date, 
				invoice_price, 
				invoice_currency, 
				invoice_country,
				tag.tag_name
			FROM
				invoice
			LEFT JOIN invoice_tag ON invoice_tag.invoice_id = invoice.invoice_id
			LEFT JOIN tag ON tag.tag_id = invoice_tag.tag_id
			WHERE
					invoice_name LIKE ?
					OR tag.tag_name LIKE ?
					OR invoice_date LIKE ?
					OR invoice_currency LIKE ?
					OR invoice_country LIKE ?
			ORDER BY
				invoice_date DESC;`

	q := "%" + c.FormValue("q") + "%"

	db := w.BillRepo.GetDb()
	rows, err := db.Query(queryBase, q, q, q, q, q)
	if err != nil {
		r["message"] = "Error while querying the database"
		return c.Render(http.StatusOK, "search-bills-result.html", r)
	}
	defer rows.Close()

	var result []BillRequest
	for rows.Next() {
		var (
			Id       string
			Name     string
			Date     string
			Price    float64
			Currency string
			Country  string
			Tag      string
		)
		rows.Scan(&Id, &Name, &Date, &Price, &Currency, &Country, &Tag)
		if err != nil {
			r["message"] = "Error while scanning the database"
			return c.Render(http.StatusOK, "search-bills-result.html", r)
		}
		result = append(result, BillRequest{
			Id:       Id,
			Name:     Name,
			Date:     Date,
			Price:    Price,
			Currency: Currency,
			// TODO exchange rate system
			ExchangeRate: 0,
			Country:      Country,
			Tag:          Tag,
		})
	}

	r["result"] = result
	r["success"] = true
	return c.Render(http.StatusOK, "search-bills-result.html", r)
}
