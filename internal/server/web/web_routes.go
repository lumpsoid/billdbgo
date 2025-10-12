package web

import (
	repository "billdb/internal/repository/bill"
	"billdb/internal/server"

	"github.com/labstack/echo/v4"
)

type WebHandlers struct {
	Config   *server.Config
	Echo     *echo.Echo
	BillRepo repository.BillRepository
}

func NewWebHandlers(config *server.Config, echo *echo.Echo, repo repository.BillRepository) *WebHandlers {
	return &WebHandlers{
		Config:   config,
		Echo:     echo,
		BillRepo: repo,
	}
}

func (w *WebHandlers) RegisterRoutes(group *echo.Group) {
	group.GET("/", w.IndexPage).Name = "index"
	group.GET("/bill/form", w.BillFormPage).Name = "bill-form"
	group.POST("/bill/form", w.BillFormSubmit)
	group.GET("/bill/link", w.BillFromLink).Name = "bill-from-link"
	group.POST("/bill/link", w.BillFromLinkResponse)
	group.GET("/bill/qr", w.BillFromQr).Name = "bill-from-qr"
	group.POST("/bill/qr", w.BillFromQrUpload)

	group.GET("/db/save", w.SaveDb).Name = "db-save"
	group.GET("/db/upload", w.UploadDb).Name = "db-upload"
	group.POST("/db/upload", w.UploadDbSubmit)

	group.GET("/browse/bills", w.BillBrowseLanding).Name = "browse-landing"
	group.GET("/browse/bills/:y/:m", w.BillBrowse).Name = "browse-bills"

	group.GET("/browse/items/:y/:m", w.ItemsBrowse).Name = "browse-items"
	group.GET("/bill/:id", w.BillView).Name = "bill-view"

	group.GET("/bill/:id/edit", w.BillEditPage).Name = "bill-edit"
	group.PUT("/bill/:id/edit", w.BillEditSubmit)

	group.GET("/search", w.SearchPage).Name = "search"
	group.GET("/search/bills", w.BillsSearch).Name = "bills-search"
	group.POST("/search/bills", w.BillSearchQueary)

	group.GET("/search/items", w.ItemsSearch).Name = "items-search"
	group.POST("/search/items", w.ItemsSearchQueary)
}
