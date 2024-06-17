package flutterapi

import "billdb/internal/server"

const baseApiPath = "/api/flutter"

type BillApi struct {
	Timestamp    int64   `json:"timestamp"`
	Name         string  `json:"name"`
	Date         string  `json:"date"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	ExchangeRate float64 `json:"exchange_rate"`
	Country      string  `json:"country"`
	Items        int     `json:"items"`
	Link         string  `json:"link"`
	Duplicates   int     `json:"duplicates"`
}

type ResponseFlutter struct {
	Success string    `json:"success"`
	Message string    `json:"message"`
	Force   bool      `json:"force"`
	Bill    []BillApi `json:"bill"`
}

func FlutterApiRoutes(s *server.Server) {
	QrHandler(s)
	FormHandler(s)
	GetTagsHandler(s)
	GetCurrenciesHandler(s)
}
