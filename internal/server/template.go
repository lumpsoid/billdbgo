package server

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		// this adding reverse function call from a template
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.Templates.ExecuteTemplate(w, name, data)
}
