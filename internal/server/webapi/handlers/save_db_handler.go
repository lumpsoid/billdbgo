package handlers

import (
	"billdb/internal/server"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

var SaveDb = server.Get("/billdb-save", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
    timestamp := time.Now().Local().Format("20060102")
    return c.Attachment(
      s.Config.Db.Path, 
      fmt.Sprintf("billdb%s.db", timestamp),
    )
	}
})
