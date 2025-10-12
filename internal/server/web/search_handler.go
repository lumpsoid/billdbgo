package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (w *WebHandlers) SearchPage(c echo.Context) error {
	return c.Render(http.StatusOK, "search.html", nil)
}
