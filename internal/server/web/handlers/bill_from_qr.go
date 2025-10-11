package handlers

import (
	"billdb/internal/parser"
	"billdb/internal/qrcode"
	"billdb/internal/server"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

var BillFromQr = server.Get("/bill/qr", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "bill-from-qr.html", map[string]interface{}{})
	}
})

var BillFromQrUpload = server.Post("/bill/qr", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		r := map[string]interface{}{
			"success": false,
		}

		err = server.CheckFormFile(file)
		if err != nil {
			r["message"] = err.Error()
			return c.Render(http.StatusOK, "bill-insert-response.html", r)
		}

		// TODO change path
		qrFilepath, err := server.UploadFileToServer(
			s.Config.QrPath,
			file,
		)
		if err != nil {
			return err
		}

		qrString, err := qrcode.ParseImage(qrFilepath)
		if err != nil {
			r["message"] = err.Error()
			r["qrPath"] = filepath.Base(qrFilepath)
			err = c.Render(http.StatusOK, "bill-insert-response.html", r)
      return err
		}
		err = os.Remove(qrFilepath)
		if err != nil {
			return err
		}

		p, err := parser.GetBillParser(qrString)
		if err != nil {
			return err
		}

		dupCheck := false
		// check for duplicates by url
		if p.Type() == "rs" {
			dupCount, err := s.BillRepo.CheckDuplicateBillByUrl(qrString)
			if err != nil {
				return err
			}
			if dupCount != 0 {
				r["success"] = false
				r["message"] = "Found duplicate bills"
				r["dupInt"] = dupCount
				return c.Render(http.StatusOK, "bill-insert-response.html", r)
			}
			dupCheck = true
		}

		b, err := p.Parse(qrString)
		if err != nil {
      r["success"] = false
      r["message"] = "Error while parsing the site"
      return c.Render(http.StatusOK, "bill-insert-response.html", r)
		}

		// if duplicates was not checked earlier
		// check it with parsed data
		if !dupCheck {
			dupCount, err := s.BillRepo.CheckDuplicateBill(b)
			if err != nil {
				return err
			}
			if dupCount != 0 {
				r["success"] = false
				r["message"] = "Found duplicate bills"
				r["dupInt"] = dupCount
				return c.Render(http.StatusOK, "bill-insert-response.html", r)
			}
			dupCheck = true
		}

		if dupCheck {
			err = s.BillRepo.InsertBillWithItems(b)
			if err != nil {
				return err
			}
		} else {
			r["success"] = false
			r["message"] = "Duplicates was not checked"
			return c.Render(http.StatusOK, "bill-insert-response.html", r)
		}

		r["success"] = true
		r["message"] = "Bill parsed successfully"
		r["bill"] = b
		return c.Render(http.StatusOK, "bill-insert-response.html", r)
	}
})
