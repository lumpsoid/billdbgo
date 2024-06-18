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
	h.SaveDb(s)
	h.BillView(s).Name = "bill-view"
	h.BillBrowse(s).Name = "bill-browse"
	h.BillEditPage(s).Name = "bill-edit"
	h.BillEditSubmit(s)
}
