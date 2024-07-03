package repository

import (
	bl "billdb/internal/bill"
	"billdb/internal/bill/item"
	"database/sql"
)

type BillRepository interface {
	GetDb() *sql.DB
	ApplyMigration(sqlFilePath string) error
	CheckDuplicateBill(bill *bl.Bill) (int, error)
	InsertBill(bill *bl.Bill) error
	GetBillByID(id string) (*bl.Bill, error)
	UpdateBill(bill *bl.Bill) error
	DeleteBill(id string) error
	InsertItems(items []*item.Item) error
	GetItemsByID(billId string) ([]*item.Item, error)
	UpdateItems(items []*item.Item) error
	DeleteItems(items []*item.Item) error
	GetCurrencies() ([]string, error)
	GetCountries() ([]string, error)
	GetTags() ([]string, error)
}
