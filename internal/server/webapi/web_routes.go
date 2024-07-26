package webapi

import (
	"billdb/internal/server"
	h "billdb/internal/server/webapi/handlers"
)

func RegisterWebRoutes(s *server.Server) {
	h.Index(s).Name = "index"
	h.BillFormPage(s).Name = "bill-form"
	h.BillFormSubmit(s)
	h.BillFromLink(s).Name = "bill-from-link"
	h.BillFromLinkResponse(s)
  h.BillFromQr(s).Name = "bill-from-qr"
  h.BillFromQrUpload(s)
	h.SaveDb(s).Name = "db-save"
	h.UploadDb(s).Name = "db-upload"
	h.UploadDbSubmit(s)
	h.BillBrowseLanding(s).Name = "browse-landing"
	h.BillBrowse(s).Name = "browse-bills"
	h.ItemsBrowse(s).Name = "browse-items"
	h.BillView(s).Name = "bill-view"
	h.BillEditPage(s).Name = "bill-edit"
	h.BillEditSubmit(s)
	h.SearchPage(s).Name = "search"
	h.BillsSearch(s).Name = "bills-search"
	h.BillSearchQueary(s)
	h.ItemsSearch(s).Name = "items-search"
	h.ItemsSearchQueary(s)
}
