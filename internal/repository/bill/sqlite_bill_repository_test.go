package repository

import (
	"billdb/internal/bill"
	"billdb/internal/bill/country"
	"billdb/internal/bill/currency"
	"billdb/internal/bill/item"
	"billdb/internal/bill/tag"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/segmentio/ksuid"
)

var dbPath string
var creationSql string

func initEnv() {
	if dbPath == "" {
		dbPath = "../../../test/db/test.db"
	}
	if creationSql == "" {
		creationSql = "./migrations/001_initial_schema.sql"
	}
}

func setUpDB(t *testing.T) (*SqliteBillRepository, error) {
	os.Remove(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Errorf("Failed to open sqlite database: %v", err)
		return nil, err
	}

	billRepository := NewSqliteBillRepository(db)
	if billRepository == nil {
		t.Errorf("Failed to create sqlite bill repository")
		return nil, err
	}
	if billRepository.DB == nil {
		t.Errorf("Failed to create sqlite database")
		return nil, err
	}
	return billRepository, nil
}

func TestCreateTables(t *testing.T) {
	// Implement tests for CreateTables function
	t.Log("Testing CreateTables function")

	initEnv()
	billRepository, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}

	err = billRepository.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	rows, err := billRepository.DB.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='invoice';")
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("Table bills was not created")
		return
	}
	rows.Close()

	rows, err = billRepository.DB.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='tag';")
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("Table bills was not created")
		return
	}
	rows.Close()
	rows, err = billRepository.DB.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='item';")
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("Table bills was not created")
		return
	}
	rows.Close()
}

func TestInsertBill(t *testing.T) {
	// Implement tests for InsertBill function
	t.Log("Testing InsertBill function")

	initEnv()
	billRepository, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}
	err = billRepository.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	date := time.Now()
	id := ksuid.New()
	b := bill.New(
		id.String(),
		"Test bill",
		date,
		100.0,
		currency.RSD,
		country.RUSSIA,
		[]*item.Item{},
		tag.New("tag1,tag2"),
		"linkString",
		"billText",
	)
	err = billRepository.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}

	var billId string
	var billName string
	var billDate string
	var billPrice float64
	var billCurrency string
	var billCountry string
	var billTag string
	var billLink string
	var billText string
	rows, err := billRepository.DB.Query(`SELECT
			invoice.invoice_id, 
			invoice_name, 
			invoice_date, 
			invoice_price, 
			invoice_currency, 
			invoice_country,
			tag.tag_name,
			invoice_link,
			invoice_text
		FROM
			invoice
		LEFT JOIN invoice_tag ON invoice_tag.invoice_id = invoice.invoice_id
		LEFT JOIN tag ON tag.tag_id = invoice_tag.tag_id
		WHERE
			invoice.invoice_id = ?`, id.String())
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&billId, &billName, &billDate, &billPrice, &billCurrency, &billCountry, &billTag, &billLink, &billText)
	if err != nil {
		t.Errorf("Failed to scan row: %v", err)
		return
	}

	if billId != b.Id {
		t.Errorf("Expected ID '%s', got %s", b.Id, billId)
	}
	if billName != b.Name {
		t.Errorf("Expected Name '%s', got %s", b.Name, billName)
	}
	if billDate != b.Date.Format("2006-01-02") {
		t.Errorf("Expected Date '%s', got %s", b.Date.String(), billDate)
	}
	if billPrice != b.Price {
		t.Errorf("Expected Price '%f', got %f", b.Price, billPrice)
	}
	if billCurrency != "rsd" {
		t.Errorf("Expected Currency '%d', got %s", b.Currency, billCurrency)
	}
	if billCountry != "russia" {
		t.Errorf("Expected Country '%d', got %s", b.Country, billCountry)
	}
	if billTag != b.Tag.String {
		t.Errorf("Expected Tag '%s', got %s", b.Tag.String, billTag)
	}
	if billLink != b.Link {
		t.Errorf("Expected Link '%s', got %s", b.Link, billLink)
	}
	if billText != b.BillText {
		t.Errorf("Expected BillText '%s', got %s", b.BillText, billText)
	}
}

func TestGetBillById(t *testing.T) {
	// Implement tests for InsertBill function
	t.Log("Testing InsertBill function")

	initEnv()
	billRepository, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}
	err = billRepository.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	date := time.Now()
	id := ksuid.New()
	b := bill.New(
		id.String(),
		"Test bill",
		date,
		100.0,
		currency.RSD,
		country.RUSSIA,
		[]*item.Item{},
		tag.New("tag1,tag2"),
		"linkString",
		"billText",
	)
	err = billRepository.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}
	billById, err := billRepository.GetBillByID(id.String())
	if err != nil {
		t.Errorf("Failed to get bill by ID: %v", err)
		return
	}

	if billById.Id != b.Id {
		t.Errorf("Expected ID '%s', got %s", b.Id, billById.Id)
	}
	if billById.Name != b.Name {
		t.Errorf("Expected Name '%s', got %s", b.Name, billById.Name)
	}
	if billById.Date.Local().Format("2006-01-02") != b.Date.Local().Format("2006-01-02") {
		t.Errorf("Expected Date '%s', got %s", b.Date.String(), billById.Date.String())
	}
	if billById.Price != b.Price {
		t.Errorf("Expected Price '%f', got %f", b.Price, billById.Price)
	}
	if billById.Currency != b.Currency {
		t.Errorf("Expected Currency '%d', got %d", b.Currency, billById.Currency)
	}
	if billById.Country != b.Country {
		t.Errorf("Expected Country '%d', got %d", b.Country, billById.Country)
	}
	if billById.Tag.String != b.Tag.String {
		t.Errorf("Expected Tag '%s', got %s", b.Tag.String, billById.Tag.String)
	}
	if billById.Link != b.Link {
		t.Errorf("Expected Link '%s', got %s", b.Link, billById.Link)
	}
	// TODO convert each type into separate package
	// and create NewWithText, New, GetBillTextOrDefault functions
	// if billById.BillText != b.BillText {
	// 	t.Errorf("Expected BillText '%s', got '%s'", b.BillText, billById.BillText)
	// }
}

func TestUpdateBill(t *testing.T) {
	// Implement tests for InsertBill function
	t.Log("Testing UpdateBill function")

	initEnv()
	billRepository, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}
	err = billRepository.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	id := ksuid.New()
	date := time.Now()
	dateNew, err := time.Parse("2006-01-02", "2023-11-13")
	if err != nil {
		t.Errorf("Failed to parse date: %v", err)
		return
	}
	bills := []*bill.Bill{
		bill.New(
			id.String(),
			"Test bill",
			date,
			100.0,
			currency.RSD,
			country.RUSSIA,
			[]*item.Item{},
			tag.Empty(),
			"linkString",
			"billText",
		),
		bill.New(
			id.String(),
			"Test bill NEW",
			dateNew,
			999.0,
			currency.EUR,
			country.SERBIA,
			[]*item.Item{},
			tag.New("tag1,tag2,NEW"),
			"linkStringNEW",
			"billTextNEW",
		),
		bill.New(
			id.String(),
			"Test bill NEW",
			dateNew,
			999.0,
			currency.EUR,
			country.SERBIA,
			[]*item.Item{},
			tag.New(""),
			"linkStringNEW",
			"billTextNEW",
		),
		bill.New(
			id.String(),
			"NEW",
			time.Now(),
			389.0,
			currency.TRY,
			country.TURKEY,
			[]*item.Item{},
			tag.Empty(),
			"linkStringNEWsadf",
			"billTextNEWxzvzcv",
		),
	}

	err = billRepository.InsertBill(bills[0])
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}

	for _, b := range bills {
		err = billRepository.UpdateBill(b)
		if err != nil {
			t.Errorf("Failed to update bill: %v", err)
			return
		}

		var billId string
		var billName string
		var billDate string
		var billPrice float64
		var billCurrency string
		var billCountry string
		var billTag *string
		var billLink string
		var billText string
		item := billRepository.DB.QueryRow(`SELECT
			invoice.invoice_id, 
			invoice_name, 
			invoice_date, 
			invoice_price, 
			invoice_currency, 
			invoice_country,
			tag.tag_name,
			invoice_link,
			invoice_text
		FROM
			invoice
		LEFT JOIN invoice_tag ON invoice_tag.invoice_id = invoice.invoice_id
		LEFT JOIN tag ON tag.tag_id = invoice_tag.tag_id
		WHERE invoice.invoice_id = ?`,
			id.String(),
		)
		err = item.Scan(&billId, &billName, &billDate, &billPrice, &billCurrency, &billCountry, &billTag, &billLink, &billText)
		if err != nil {
			t.Errorf("Failed to scan row: %v", err)
			return
		}

		if b.Id != billId {
			t.Errorf("Expected ID '%s', got %s", b.Id, billId)
		}
		if b.Name != billName {
			t.Errorf("Expected Name '%s', got %s", b.Name, billName)
		}
		if b.Date.Format("2006-01-02") != billDate {
			t.Errorf(
				"Expected Date '%s', got '%s'",
				b.Date.Format("2006-01-02"),
				billDate,
			)
		}
		if b.Price != billPrice {
			t.Errorf("Expected Price '%f', got %f", b.Price, billPrice)
		}
		if b.Currency.String() != billCurrency {
			t.Errorf("Expected Currency '%s', got %s", b.Currency, billCurrency)
		}
		if b.Country.String() != billCountry {
			t.Errorf("Expected Country '%s', got %s", b.Country, billCountry)
		}
		if b.Tag.Valid && billTag == nil {
			t.Errorf(
				"Expected Tag to be valid: %v, got %v\n",
				b.Tag.Valid,
				billTag == nil,
			)
		}
		if !b.Tag.Valid == (billTag != nil) {
			t.Errorf(
				"Expected Tag to be not valid: %v, got %v\n",
				b.Tag.Valid,
				billTag != nil,
			)
		}
		if b.Tag.Valid {
			if b.Tag.String != *billTag {
				t.Errorf("Expected Tag '%s', got '%s'", b.Tag.String, *billTag)
			}
		}
		if b.Link != billLink {
			t.Errorf("Expected Link '%s', got %s", b.Link, billLink)
		}
		if b.BillText != billText {
			t.Errorf("Expected BillText '%s', got %s", b.BillText, billText)
		}
	}
}

func TestDeleteBill(t *testing.T) {
	t.Log("Testing DeleteBill function")

	initEnv()
	billRepository, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}
	err = billRepository.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	date := time.Now()
	id := ksuid.New()
	b := bill.New(
		id.String(),
		"Test bill",
		date,
		100.0,
		currency.RSD,
		country.RUSSIA,
		[]*item.Item{},
		tag.New("tag1,tag2"),
		"linkString",
		"billText",
	)
	err = billRepository.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}

	err = billRepository.DeleteBill(b.Id)
	if err != nil {
		t.Errorf("Failed to delete bill: %v", err)
		return
	}

	query :=
		`SELECT * FROM invoice WHERE invoice_id = ?`

	rows, err := billRepository.DB.Query(query, b.Id)
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		t.Errorf("Failed to delete bill: %v", err)
		return
	}

	query =
		`SELECT * FROM invoice_tag WHERE invoice_id = ?`

	rows, err = billRepository.DB.Query(query, b.Id)
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		t.Errorf("Failed to delete bill: %v", err)
		return
	}
}

func TestCheckDuplicateBill(t *testing.T) {
	t.Log("Testing DuplicatesBills function")

	initEnv()
	billRepository, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}
	err = billRepository.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	date := time.Now()
	id := ksuid.New()
	b := bill.New(
		id.String(),
		"Test bill",
		date,
		100.0,
		currency.RSD,
		country.RUSSIA,
		[]*item.Item{},
		tag.New("tag1,tag2"),
		"linkString",
		"billText",
	)
	err = billRepository.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}

	billDupCount, err := billRepository.CheckDuplicateBill(b)
	if err != nil {
		t.Errorf("Failed to get duplicate bills: %v", err)
		return
	}

	if billDupCount != 1 {
		t.Errorf("Expected 1 duplicate bill, got %d", billDupCount)
		return
	}

}

func TestScanToBill(t *testing.T) {
	fmt.Println("Testing ScanBill function")
	initEnv()
	billRepo, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}
	err = billRepo.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}
	date := time.Now()
	id := ksuid.New()
	b := bill.New(
		id.String(),
		"Test bill",
		date,
		100.0,
		currency.RSD,
		country.RUSSIA,
		[]*item.Item{},
		tag.New("tag1,tag2"),
		"linkString",
		"billText",
	)
	err = billRepo.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}

	query := `SELECT
			invoice.invoice_id, 
			invoice_name, 
			invoice_date, 
			invoice_price, 
			invoice_currency, 
			invoice_country,
			tag.tag_name,
			invoice_link
		FROM invoice
		LEFT JOIN invoice_tag ON invoice_tag.invoice_id = invoice.invoice_id
		LEFT JOIN tag ON tag.tag_id = invoice_tag.tag_id
		WHERE invoice.invoice_id = ?`
	row := billRepo.DB.QueryRow(query, b.Id)
	bN, err := ScanToBill(row)
	if err != nil {
		t.Errorf("Failed to scan bill: %v", err)
		return
	}
	if bN.Id != b.Id {
		t.Errorf("Expected ID '%s', got %s", b.Id, bN.Id)
	}
	if bN.Name != b.Name {
		t.Errorf("Expected Name '%s', got %s", b.Name, bN.Name)
	}
	if bN.Date.Format("2006-01-02") != b.Date.Format("2006-01-02") {
		t.Errorf("Expected Date '%s', got %s", b.Date.String(), bN.Date.String())
	}
	if bN.Price != b.Price {
		t.Errorf("Expected Price '%f', got %f", b.Price, bN.Price)
	}
	if bN.Currency != b.Currency {
		t.Errorf("Expected Currency '%d', got %d", b.Currency, bN.Currency)
	}
	if bN.Country != b.Country {
		t.Errorf("Expected Country '%d', got %d", b.Country, bN.Country)
	}
	if bN.Tag.String != b.Tag.String {
		t.Errorf("Expected Tag '%s', got %s", b.Tag.String, bN.Tag.String)
	}
	if bN.Link != b.Link {
		t.Errorf("Expected Link '%s', got %s", b.Link, bN.Link)
	}
}

func TestInsertItems(t *testing.T) {
	fmt.Println("Testing InsertItems function")
	initEnv()
	billRepo, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}
	err = billRepo.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}
	date := time.Now()
	id := ksuid.New()
	b := bill.New(
		id.String(),
		"Test bill",
		date,
		100.0,
		currency.RSD,
		country.RUSSIA,
		[]*item.Item{},
		tag.New("tag1,tag2"),
		"linkString",
		"billText",
	)

	itemId := ksuid.New()
	var items []*item.Item
	itemN := item.New(
		itemId.String(),
		id.String(),
		"item1",
		100.0,
		100.0,
		1.0,
	)
	items = append(items, itemN)
	err = billRepo.InsertItems(items)
	if err != nil {
		t.Errorf("Failed to insert items: %v", err)
		return
	}

	rows, err := billRepo.DB.Query(`SELECT
		item_id, 
		invoice_id, 
		item_name, 
		item_price, 
		item_price_one, 
		item_quantity
	FROM item 
	WHERE invoice_id = ?`, b.Id)
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var itemId string
		var billId string
		var name string
		var price float64
		var priceOne float64
		var quantity float64
		err = rows.Scan(&itemId, &billId, &name, &price, &priceOne, &quantity)
		if err != nil {
			t.Errorf("Failed to scan row: %v", err)
			return
		}
		if itemId != itemN.ItemId {
			t.Errorf("Expected ID: '%s', got: %s", itemN.ItemId, itemId)
		}
		if billId != itemN.BillId {
			t.Errorf("Expected ID: '%s', got: %s", itemN.BillId, billId)
		}
		if name != itemN.Name {
			t.Errorf("Expected Name '%s', got %s", itemN.Name, name)
		}
		if price != itemN.Price {
			t.Errorf("Expected Price '%f', got %f", itemN.Price, price)
		}
		if priceOne != itemN.PriceOne {
			t.Errorf("Expected PriceOne '%f', got %f", itemN.PriceOne, priceOne)
		}
		if quantity != itemN.Quantity {
			t.Errorf("Expected Quantity '%f', got %f", itemN.Quantity, quantity)
		}
	}
}

func TestGetItemsByID(t *testing.T) {
	fmt.Println("Testing GetItemsByID function")
	initEnv()
	billRepo, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}
	err = billRepo.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}
	id := ksuid.New()
	itemId1 := ksuid.New()
	itemN := item.New(
		itemId1.String(),
		id.String(),
		"item1",
		100.0,
		100.0,
		1.0,
	)
	itemId2 := ksuid.New()
	item2 := item.New(
		itemId2.String(),
		id.String(),
		"item2",
		102.0,
		102.0,
		2.0,
	)
	err = billRepo.InsertItems([]*item.Item{
		itemN,
		item2,
	})
	if err != nil {
		t.Errorf("Failed to insert items: %v", err)
		return
	}

	itemsByID, err := billRepo.GetItemsByID(itemN.BillId)
	if err != nil {
		t.Errorf("Failed to get items by ID: %v", err)
		return
	}

	if len(itemsByID) != 2 {
		t.Errorf("Expected 2 item, got %d", len(itemsByID))
		return
	}
	if itemsByID[0].BillId != itemN.BillId {
		t.Errorf("Expected ID '%s', got %s", itemN.BillId, itemsByID[0].BillId)
	}
	if itemsByID[0].Name != itemN.Name {
		t.Errorf("Expected Name '%s', got %s", itemN.Name, itemsByID[0].Name)
	}
	if itemsByID[1].ItemId != item2.ItemId {
		t.Errorf("Expected Price '%f', got %f", item2.Price, itemsByID[1].Price)
	}
	if itemsByID[1].BillId != item2.BillId {
		t.Errorf("Expected Price '%f', got %f", item2.Price, itemsByID[1].Price)
	}
	if itemsByID[1].Price != item2.Price {
		t.Errorf("Expected Price '%f', got %f", item2.Price, itemsByID[1].Price)
	}
}

func TestDeleteItems(t *testing.T) {
	t.Log("Testing DeleteItems function")
	initEnv()
	billRepo, err := setUpDB(t)
	if err != nil {
		t.Errorf("Failed to set up database: %v", err)
		return
	}
	err = billRepo.ApplyMigration(creationSql)
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}
	id := ksuid.New()
	itemId1 := ksuid.New()
	var items []*item.Item
	itemN := item.New(
		itemId1.String(),
		id.String(),
		"item1",
		100.0,
		100.0,
		1.0,
	)
	itemId2 := ksuid.New()
	item2 := item.New(
		itemId2.String(),
		id.String(),
		"item2",
		102.0,
		102.0,
		2.0,
	)
	items = append(items, itemN, item2)
	err = billRepo.InsertItems(items)
	if err != nil {
		t.Errorf("Failed to insert items: %v", err)
		return
	}

	err = billRepo.DeleteItems(items)
	if err != nil {
		t.Errorf("Failed to delete items by ID: %v", err)
		return
	}

	itemsFromDb, err := billRepo.GetItemsByID(itemN.BillId)
	if err != nil {
		t.Errorf("Failed to get items by ID: %v", err)
		return
	}

	if len(itemsFromDb) != 0 {
		t.Errorf("Expected 0 item, got %d", len(itemsFromDb))
		return
	}

	rows, err := billRepo.DB.Query(`SELECT item_id
		FROM item_tag 
		WHERE item_id = ?`,
		itemN.ItemId,
	)
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()
	if rows.Next() {
		t.Error("Expected no rows in item_tag table")
	}
	rows.Close()
}
