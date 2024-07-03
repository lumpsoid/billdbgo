package repository

import (
	bl "billdb/internal/bill"
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
	InsertItems(items []*bl.Item) error
	GetItemsByID(billId string) ([]*bl.Item, error)
	UpdateItems(items []*bl.Item) error
	DeleteItems(items []*bl.Item) error
	GetCurrencies() ([]string, error)
	GetCountries() ([]string, error)
	GetTags() ([]string, error)
}
