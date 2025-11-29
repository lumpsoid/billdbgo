package web

import (
	"billdb/internal/parser"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool
	Message string
	Bill    BillRequest
}

func (w *WebHandlers) BillFromLink(c echo.Context) error {
	return c.Render(http.StatusOK, "bill-from-link.html", map[string]any{})
}

func (w *WebHandlers) BillFromLinkResponse(c echo.Context) error {
	r := map[string]any{
		"success": false,
		"results": []map[string]any{},
	}
	linkText := c.FormValue("link")
	// Split by newlines and process each link
	lines := strings.Split(linkText, "\n")
	var validLinks []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			validLinks = append(validLinks, trimmed)
		}
	}
	if len(validLinks) == 0 {
		r["message"] = "No valid links provided"
		return c.Render(http.StatusOK, "bill-insert-response.html", r)
	}
	// Process each valid link
	successCount := 0
	for _, link := range validLinks {
		// Truncate link for display
		linkDisplay := link
		if len(link) > 10 {
			linkDisplay = "..." + link[len(link)-10:]
		}

		linkResult := map[string]any{
			"link":        link,
			"linkDisplay": linkDisplay,
			"success":     false,
		}
		p, err := parser.GetBillParser(link)
		if err != nil {
			linkResult["message"] = err.Error()
			r["results"] = append(r["results"].([]map[string]any), linkResult)
			continue
		}
		dupCheck := false
		if p.Type() == "rs" {
			dupCount, err := w.BillRepo.CheckDuplicateBillByUrl(link)
			if err != nil {
				linkResult["message"] = err.Error()
				r["results"] = append(r["results"].([]map[string]any), linkResult)
				continue
			}
			if dupCount != 0 {
				linkResult["message"] = fmt.Sprintf("Found %d duplicate bills", dupCount)
				linkResult["dupInt"] = dupCount
				r["results"] = append(r["results"].([]map[string]any), linkResult)
				continue
			}
			dupCheck = true
		}
		b, err := p.Parse(link)
		if err != nil {
			linkResult["message"] = "Error while parsing the site"
			r["results"] = append(r["results"].([]map[string]any), linkResult)
			continue
		}
		if dupCheck {
			err = w.BillRepo.InsertBillWithItems(b)
			if err != nil {
				linkResult["message"] = err.Error()
				r["results"] = append(r["results"].([]map[string]any), linkResult)
				continue
			}
		} else {
			linkResult["message"] = "Duplicates was not checked"
			r["results"] = append(r["results"].([]map[string]any), linkResult)
			continue
		}
		linkResult["success"] = true
		linkResult["message"] = "Bill parsed successfully"
		linkResult["bill"] = b
		successCount++
		r["results"] = append(r["results"].([]map[string]any), linkResult)
	}
	if successCount > 0 {
		r["success"] = true
	}
	r["message"] = fmt.Sprintf("Processed %d links: %d successful, %d failed",
		len(validLinks), successCount, len(validLinks)-successCount)
	return c.Render(http.StatusOK, "bill-insert-response.html", r)
}
