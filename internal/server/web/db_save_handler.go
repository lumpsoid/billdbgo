package web

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

func (w *WebHandlers) SaveDb(c echo.Context) error {
	timestamp := time.Now().Local().Format("20060102")
	return c.Attachment(
		w.Config.DbPath,
		fmt.Sprintf("billdb%s.db", timestamp),
	)
}
