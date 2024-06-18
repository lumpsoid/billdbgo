package repository

import (
	bl "billdb/internal/bill"
	"database/sql"
)

type BillRepository interface {
	GetDb() *sql.DB
	CreateTables() error
	CheckUniqueItemNames() error
	CheckDuplicateBill(bill *bl.Bill) ([]*bl.Bill, error)
	InsertBill(bill *bl.Bill) error
	GetBillByID(id int64) (*bl.Bill, error)
	UpdateBill(bill *bl.Bill) error
	DeleteBill(id int64) error
	InsertItems(items []*bl.Item) error
	GetItemsByID(billId int64) ([]*bl.Item, error)
	UpdateItems(items []*bl.Item) error
	DeleteItems(billId int64) error
	GetCurrencies() ([]string, error)
	GetCountries() ([]string, error)
	GetTags() ([]string, error)
}
