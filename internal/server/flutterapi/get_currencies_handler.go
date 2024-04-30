package flutterapi

import (
	"billdb/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

var GetCurrenciesHandler = server.Get(baseApiPath+"/currencies", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := s.BillRepo.GetDb()
		rows, err := db.Query("SELECT DISTINCT currency FROM bills;")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		defer rows.Close()
		var currencies []string
		for rows.Next() {
			var currency string
			err := rows.Scan(&currency)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			currencies = append(currencies, currency)
		}
		return c.JSON(http.StatusOK, currencies)
	}
})
