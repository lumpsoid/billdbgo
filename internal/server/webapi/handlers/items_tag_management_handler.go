package handlers

import (
	"billdb/internal/server"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ItemsTagManagement = server.Get("/tags/items", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			htmlPath := "tags-items.html"
			r := make(map[string]interface{})
			r["Success"] = false

			query := `SELECT
				m.id, m.name, m.tag
			FROM
				items_meta_v2 as m
			LEFT JOIN items as i ON m.name = i.name
			LEFT JOIN bills as b ON i.id = b.id
			ORDER BY
				b.dates DESC;`

			db := s.BillRepo.GetDb()
			rows, err := db.Query(query)
			if err != nil {
				r["Message"] = "Error while querying the database"
				return c.Render(http.StatusOK, htmlPath, r)
			}
			defer rows.Close()

			var itemsResponse []map[string]interface{}
			for rows.Next() {
				var (
					Id   int64
					Name string
					Tag  string
				)
				rows.Scan(&Id, &Name, &Tag)
				if err != nil {
					r["Message"] = "Error while scanning the database"
					return c.Render(http.StatusOK, htmlPath, r)
				}
				itemsResponse = append(itemsResponse, map[string]interface{}{
					"Id":   Id,
					"Name": Name,
					"Tag":  Tag,
				})
			}

			r["Items"] = itemsResponse
			r["Success"] = true
			return c.Render(http.StatusOK, htmlPath, r)
		}
	})

	ItemTagGet = server.Get("/tag/item/:id", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			htmlPath := "tags-items-info.html"
			r := make(map[string]interface{})
			r["Success"] = false

			id := c.Param("id")
			query := `SELECT
				id, name, tag
			FROM
				items_meta_v2
			WHERE
				id = ?;`

			db := s.BillRepo.GetDb()
			row := db.QueryRow(query, id)

			var (
				Id   int64
				Name string
				Tag  string
			)
			err := row.Scan(&Id, &Name, &Tag)
			if err != nil {
				if err == sql.ErrNoRows {
					r["Message"] = "No such item found in the database"
					return c.Render(http.StatusOK, htmlPath, r)
				}
				r["Message"] = "Error while querying the database"
				return c.Render(http.StatusOK, htmlPath, r)
			}

			r["Id"] = Id
			r["Name"] = Name
			r["Tag"] = Tag
			r["Success"] = true
			return c.Render(http.StatusOK, htmlPath, r)
		}
	})

	ItemTagEdit = server.Get("/tag/item/:id/edit", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			htmlPath := "tags-items-edit.html"
			r := make(map[string]interface{})
			r["Success"] = false

			id := c.Param("id")
			query := `SELECT
				id, name, tag
			FROM
				items_meta_v2
			WHERE
				id = ?;`

			db := s.BillRepo.GetDb()
			row := db.QueryRow(query, id)

			var (
				Id   int64
				Name string
				Tag  string
			)
			err := row.Scan(&Id, &Name, &Tag)
			if err != nil {
				if err == sql.ErrNoRows {
					r["Message"] = "No such item found in the database"
					return c.Render(http.StatusOK, htmlPath, r)
				}
				r["Message"] = "Error while querying the database"
				return c.Render(http.StatusOK, htmlPath, r)
			}

			r["Id"] = Id
			r["Name"] = Name
			r["Tag"] = Tag
			r["Success"] = true
			return c.Render(http.StatusOK, htmlPath, r)
		}
	})

	ItemTagUpdate = server.Put("/tag/item/:id", func(s *server.Server) echo.HandlerFunc {
		return func(c echo.Context) error {
			htmlPath := "tags-items-result.html"
			r := make(map[string]interface{})
			r["Success"] = false

			id := c.Param("id")
			tagNew := c.FormValue("tag")
			r["Id"] = id

			db := s.BillRepo.GetDb()

			query := `SELECT tag FROM items_meta_v2 WHERE id = ?;`
			row := db.QueryRow(query, id)
			var TagPrevious string
			err := row.Scan(&TagPrevious)
			if err != nil {
				if err == sql.ErrNoRows {
					r["Message"] = "Error while querying the database"
					return c.Render(http.StatusOK, htmlPath, r)
				}
				r["Message"] = "Error while querying the database"
				return c.Render(http.StatusOK, htmlPath, r)
			}
			r["TagPrevious"] = TagPrevious

			query = `UPDATE items_meta_v2 SET tag = ?  WHERE id = ?;`
			_, err = db.Exec(query, tagNew, id)
			if err != nil {
				r["Message"] = "Error while querying the database"
				return c.Render(http.StatusOK, htmlPath, r)
			}

			query = `SELECT id, name, tag FROM items_meta_v2 WHERE id = ?;`
			row = db.QueryRow(query, id)
			var (
				Id   int64
				Name string
				Tag  string
			)
			err = row.Scan(&Id, &Name, &Tag)
			if err != nil {
				if err == sql.ErrNoRows {
					r["Message"] = "Error while querying the database"
					return c.Render(http.StatusOK, htmlPath, r)
				}
				r["Message"] = "Error while querying the database"
				return c.Render(http.StatusOK, htmlPath, r)
			}

			r["Id"] = Id
			r["Name"] = Name
			r["Tag"] = Tag
			r["Success"] = true
			return c.Render(http.StatusOK, htmlPath, r)
		}
	})
)
