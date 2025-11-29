package web

import (
	"fmt"
	"path"
	"time"

	"github.com/labstack/echo/v4"
)

func (w *WebHandlers) SaveDb(c echo.Context) error {
	filenameTemplate := w.Config.DbFileNameTemplate
	if filenameTemplate == "" {
		filenameTemplate = "billdb-20060102"
	}

	basename := time.Now().Local().Format(filenameTemplate)
	filename := basename + ".db"
	filename = path.Base(filename)

	// Prevent caching
	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")

	c.Response().Header().Set("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, filename))

	return c.File(w.Config.DbPath)
}
