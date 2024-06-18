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
	h.SaveDb(s).Name = "db-save"
	h.UploadDb(s).Name = "db-upload"
	h.UploadDbSubmit(s)
	h.BillView(s).Name = "bill-view"
	h.BillBrowse(s).Name = "bills"
	h.BillEditPage(s).Name = "bill-edit"
	h.BillEditSubmit(s)
}
