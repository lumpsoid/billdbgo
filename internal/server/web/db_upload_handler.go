package web

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func (w *WebHandlers) UploadDb(c echo.Context) error {
	return c.Render(http.StatusOK, "db-upload.html", nil)
}

func (w *WebHandlers) UploadDbSubmit(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	tmpDir := os.TempDir()
	dstPath := filepath.Join(tmpDir, file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	src.Close()
	dst.Close()

	dbCurrent, err := os.Create(w.Config.DbPath)
	if err != nil {
		return err
	}

	dbNew, err := os.Open(dstPath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(dbCurrent, dbNew); err != nil {
		return err
	}

	err = os.Remove(dstPath)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/")
}
