package handlers

import (
	"billdb/internal/server"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

var UploadDb = server.Get("/db/upload", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "db-upload.html", nil)
	}
})

var UploadDbSubmit = server.Post("/db/upload", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		dbCurrent, err := os.Create(s.Config.DbPath)
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
})
