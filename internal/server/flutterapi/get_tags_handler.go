package flutterapi

import (
	"billdb/internal/server"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var GetTagsHandler = server.Get(baseApiPath+"/tags", func(s *server.Server) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := s.BillRepo.GetDb()
		rows, err := db.Query(`
    SELECT DISTINCT tag 
    FROM bills 
    WHERE dates > '2023-07-01' 
    AND tag IS NOT NULL 
    AND tag <> ''
    UNION
    SELECT DISTINCT tag 
    FROM items_meta 
    WHERE tag IS NOT NULL 
    AND tag <> '';
    `)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		}
		defer rows.Close()

		var tags []string
		for rows.Next() {
			var tag string
			err := rows.Scan(&tag)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
			}
			tags = append(tags, tag)
		}
		return c.JSON(http.StatusOK, tags)
	}
})
