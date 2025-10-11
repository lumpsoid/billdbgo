package handlers

import (
	"billdb/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	SearchPage = server.Get("/search", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			return c.Render(http.StatusOK, "search.html", map[string]interface{}{})
		}
	})
)
