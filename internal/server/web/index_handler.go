package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (w *WebHandlers) IndexPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", map[string]any{
		"version": "0.6.2",
	})
}
