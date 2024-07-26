package handlers

import (
	"billdb/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

var Index = server.Get("/", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"version": "0.6.0",
		})
	}
})
