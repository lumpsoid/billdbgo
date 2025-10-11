package handlers

import (
	"billdb/internal/parser"
	"billdb/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool
	Message string
	Bill    BillRequest
}

var BillFromLink = server.Get("/bill/link", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.Render(http.StatusOK, "bill-from-link.html", map[string]interface{}{})
	}
})

var BillFromLinkResponse = server.Post("/bill/link", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := map[string]interface{}{
			"success": false,
		}
		link := c.FormValue("link")

		p, err := parser.GetBillParser(link)
		if err != nil {
      r["message"] = err.Error()
			return c.Render(http.StatusOK, "bill-insert-response.html", r)
		}

    dupCheck := false
    if p.Type() == "rs" {
      dupCount, err := s.BillRepo.CheckDuplicateBillByUrl(link)
      if err != nil {
        return err
      }
      if dupCount != 0 {
        r["message"] = "Found duplicate bills"
        r["dupInt"] = dupCount
        return c.Render(http.StatusOK, "bill-insert-response.html", r)
      }
      dupCheck = true
    }

		b, err := p.Parse(link)
		if err != nil {
      r["message"] = "Error while parsing the site"
			return c.Render(http.StatusOK, "bill-insert-response.html", r)
		}

    if dupCheck {
      err = s.BillRepo.InsertBillWithItems(b)
      if err != nil {
        r["message"] = err.Error()
        return c.Render(http.StatusOK, "bill-insert-response.html", r)
      }
    } else {
        r["message"] = "Duplicates was not checked"
        return c.Render(http.StatusOK, "bill-insert-response.html", r)
    }

    r["success"] = true
    r["message"] = "Bill parsed successfully"
		r["bill"] = b
		return c.Render(http.StatusOK, "bill-insert-response.html", r)
	}
},
)
