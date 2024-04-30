package repository

import (
	"billdb/internal/bill"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
)

func setUpDB(t *testing.T) (*SqliteBillRepository, error) {
	os.Remove("./test.db")

	db, err := sql.Open("sqlite3", "./test.db")
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
	fmt.Println("Testing CreateTables function")

  billRepository, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }

	err = billRepository.CreateTables()
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	rows, err := billRepository.DB.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='bills';")
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("Table bills was not created")
		return
	}
}

func TestInsertBill(t *testing.T) {
	// Implement tests for InsertBill function
	fmt.Println("Testing InsertBill function")

  billRepository, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
	err = billRepository.CreateTables()
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	timestamp := time.Now()
	b := bill.BillNew(
		timestamp,
		"Test bill",
		timestamp,
		100.0,
		bill.RSD,
		1.0,
		bill.RUSSIA,
		[]*bill.Item{},
		"tag1,tag2",
		"linkString",
		"billText",
	)
	err = billRepository.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}

	var billId int64
	var billName string
	var billDate string
	var billPrice float64
	var billCurrency string
	var billExchangeRate float64
	var billCountry string
	var billTag string
	var billLink string
	var billText string
	rows, err := billRepository.DB.Query(`
	SELECT 
		id, 
		name, 
		dates, 
		price,
		currency,
		exchange_rate,
		country,
		tag,
		link,
		bill
	FROM bills
	LIMIT 1;
	`)
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&billId, &billName, &billDate, &billPrice, &billCurrency, &billExchangeRate, &billCountry, &billTag, &billLink, &billText)
	if err != nil {
		t.Errorf("Failed to scan row: %v", err)
		return
	}

	if billId != b.Id.Local().Local().UnixMilli() {
		t.Errorf("Expected ID '%d', got %d", b.Id.Local().Local().UnixMilli(), billId)
	}
	if billName != b.Name {
		t.Errorf("Expected Name '%s', got %s", b.Name, billName)
	}
	if billDate != b.Id.Local().Format("2006-01-02") {
		t.Errorf("Expected Date '%s', got %s", b.Date.String(), billDate)
	}
	if billPrice != b.Price {
		t.Errorf("Expected Price '%f', got %f", b.Price, billPrice)
	}
	if billCurrency != "rsd" {
		t.Errorf("Expected Currency '%d', got %s", b.Currency, billCurrency)
	}
	if billExchangeRate != b.ExchangeRate {
		t.Errorf("Expected ExchangeRate '%f', got %f", b.ExchangeRate, billExchangeRate)
	}
	if billCountry != "russia" {
		t.Errorf("Expected Country '%d', got %s", b.Country, billCountry)
	}
	if billTag != string(b.Tag) {
		t.Errorf("Expected Tag '%s', got %s", b.Tag, billTag)
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
	fmt.Println("Testing InsertBill function")

  billRepository, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
	err = billRepository.CreateTables()
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	timestamp := time.Now()
	b := bill.BillNew(
		timestamp,
		"Test bill",
		timestamp,
		100.0,
		bill.RSD,
		1.0,
		bill.RUSSIA,
		[]*bill.Item{},
		"tag1,tag2",
		"linkString",
		"billText",
	)
	err = billRepository.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}
	billById, err := billRepository.GetBillByID(timestamp.UnixMilli())
	if err != nil {
		t.Errorf("Failed to get bill by ID: %v", err)
		return
	}

	if billById.Id.Local().UnixMilli() != b.Id.Local().Local().UnixMilli() {
		t.Errorf("Expected ID '%d', got %d", b.Id.Local().Local().UnixMilli(), billById.Id.Local().UnixMilli())
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
	if billById.ExchangeRate != b.ExchangeRate {
		t.Errorf("Expected ExchangeRate '%f', got %f", b.ExchangeRate, billById.ExchangeRate)
	}
	if billById.Country != b.Country {
		t.Errorf("Expected Country '%d', got %d", b.Country, billById.Country)
	}
	if billById.Tag != b.Tag {
		t.Errorf("Expected Tag '%s', got %s", b.Tag, billById.Tag)
	}
	if billById.Link != b.Link {
		t.Errorf("Expected Link '%s', got %s", b.Link, billById.Link)
	}
	if billById.BillText != b.BillText {
		t.Errorf("Expected BillText '%s', got %s", b.BillText, billById.BillText)
	}
}

func TestUpdateBill(t *testing.T) {
	// Implement tests for InsertBill function
	fmt.Println("Testing UpdateBill function")

  billRepository, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
	err = billRepository.CreateTables()
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	timestamp := time.Now()
	b := bill.BillNew(
		timestamp,
		"Test bill",
		timestamp,
		100.0,
		bill.RSD,
		1.0,
		bill.RUSSIA,
		[]*bill.Item{},
		"tag1,tag2",
		"linkString",
		"billText",
	)
	err = billRepository.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}

	timestampNew := time.Now()
	bN := bill.BillNew(
		timestampNew,
		"Test bill NEW",
		timestampNew,
		999.0,
		bill.EUR,
		120.0,
		bill.SERBIA,
		[]*bill.Item{},
		"tag1,tag2,NEW",
		"linkStringNEW",
		"billTextNEW",
	)
  err = billRepository.UpdateBill(bN)
	if err != nil {
		t.Errorf("Failed to update bill: %v", err)
		return
	}

	if bN.Id.Local().UnixMilli() == b.Id.Local().Local().UnixMilli() {
		t.Errorf("Expected ID '%d', got %d", b.Id.Local().Local().UnixMilli(), bN.Id.Local().UnixMilli())
	}
	if bN.Name == b.Name {
		t.Errorf("Expected Name '%s', got %s", b.Name, bN.Name)
	}
	if bN.Date.Local() == b.Date.Local() {
		t.Errorf("Expected Date '%s', got %s", b.Date.String(), bN.Date.String())
	}
	if bN.Price == b.Price {
		t.Errorf("Expected Price '%f', got %f", b.Price, bN.Price)
	}
	if bN.Currency == b.Currency {
		t.Errorf("Expected Currency '%d', got %d", b.Currency, bN.Currency)
	}
	if bN.ExchangeRate == b.ExchangeRate {
		t.Errorf("Expected ExchangeRate '%f', got %f", b.ExchangeRate, bN.ExchangeRate)
	}
	if bN.Country == b.Country {
		t.Errorf("Expected Country '%d', got %d", b.Country, bN.Country)
	}
	if bN.Tag == b.Tag {
		t.Errorf("Expected Tag '%s', got %s", b.Tag, bN.Tag)
	}
	if bN.Link == b.Link {
		t.Errorf("Expected Link '%s', got %s", b.Link, bN.Link)
	}
	if bN.BillText == b.BillText {
		t.Errorf("Expected BillText '%s', got %s", b.BillText, bN.BillText)
	}
}

func TestDeleteBill(t *testing.T) {
	fmt.Println("Testing DeleteBill function")

  billRepository, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
	err = billRepository.CreateTables()
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	timestamp := time.Now()
	b := bill.BillNew(
		timestamp,
		"Test bill",
		timestamp,
		100.0,
		bill.RSD,
		1.0,
		bill.RUSSIA,
		[]*bill.Item{},
		"tag1,tag2",
		"linkString",
		"billText",
	)
	err = billRepository.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}

  err = billRepository.DeleteBill(b.Id.Local().UnixMilli())
	if err != nil {
		t.Errorf("Failed to delete bill: %v", err)
		return
	}

  rows, err := billRepository.DB.Query("SELECT * FROM bills WHERE id = ?", b.Id.Local().UnixMilli())
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
	fmt.Println("Testing DuplicatesBills function")

  billRepository, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
	err = billRepository.CreateTables()
	if err != nil {
		t.Errorf("Failed to create tables: %v", err)
		return
	}

	timestamp := time.Now()
	b := bill.BillNew(
		timestamp,
		"Test bill",
		timestamp,
		100.0,
		bill.RSD,
		1.0,
		bill.RUSSIA,
		[]*bill.Item{},
		"tag1,tag2",
		"linkString",
		"billText",
	)
	err = billRepository.InsertBill(b)
	if err != nil {
		t.Errorf("Failed to insert bill: %v", err)
		return
	}

  var bills []*bill.Bill
  bills, err = billRepository.CheckDuplicateBill(b)
  if err != nil {
    t.Errorf("Failed to get duplicate bills: %v", err)
    return
  }

  if len(bills) != 1 {
    t.Errorf("Expected 1 duplicate bill, got %d", len(bills))
    return
  }

}

func TestScanToBill(t *testing.T) {
  fmt.Println("Testing ScanBill function")
  billRepo, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
  err = billRepo.CreateTables()
  if err != nil {
    t.Errorf("Failed to create tables: %v", err)
    return
  }
	timestamp := time.Now()
	b := bill.BillNew(
		timestamp,
		"Test bill",
		timestamp,
		100.0,
		bill.RSD,
		1.0,
		bill.RUSSIA,
		[]*bill.Item{},
		"tag1,tag2",
		"linkString",
		"billText",
	)
  err = billRepo.InsertBill(b)

  row := billRepo.DB.QueryRow(`SELECT
    id,
    name,
    dates,
    price,
    currency,
    exchange_rate,
    country,
    tag,
    link,
    bill
    FROM bills WHERE id = ?`, b.Id.Local().UnixMilli())
  bN, err := scanToBill(row)
  if err != nil {
    t.Errorf("Failed to scan bill: %v", err)
    return
  }
	if bN.Id.Local().UnixMilli() != b.Id.Local().Local().UnixMilli() {
		t.Errorf("Expected ID '%d', got %d", b.Id.Local().Local().UnixMilli(), bN.Id.Local().UnixMilli())
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
	if bN.ExchangeRate != b.ExchangeRate {
		t.Errorf("Expected ExchangeRate '%f', got %f", b.ExchangeRate, bN.ExchangeRate)
	}
	if bN.Country != b.Country {
		t.Errorf("Expected Country '%d', got %d", b.Country, bN.Country)
	}
	if bN.Tag != b.Tag {
		t.Errorf("Expected Tag '%s', got %s", b.Tag, bN.Tag)
	}
	if bN.Link != b.Link {
		t.Errorf("Expected Link '%s', got %s", b.Link, bN.Link)
	}
	if bN.BillText != b.BillText {
		t.Errorf("Expected BillText '%s', got %s", b.BillText, bN.BillText)
	}
}

func TestInsertItems(t *testing.T) {
  fmt.Println("Testing InsertItems function")
  billRepo, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
  err = billRepo.CreateTables()
  if err != nil {
    t.Errorf("Failed to create tables: %v", err)
    return
  }
  timestamp := time.Now()
  b := bill.BillNew(
    timestamp,
    "Test bill",
    timestamp,
    100.0,
    bill.RSD,
    1.0,
    bill.RUSSIA,
    []*bill.Item{},
    "tag1,tag2",
    "linkString",
    "billText",
  )

  var items []*bill.Item
  item := bill.ItemNew(
    time.Now(),
    "item1",
    100.0,
    100.0,
    1.0,
  )
  items = append(items, item)
  err = billRepo.InsertItems(items)
  if err != nil {
    t.Errorf("Failed to insert items: %v", err)
    return
  }

  rows, err := billRepo.DB.Query("SELECT id, name, price, price_one, quantity FROM items WHERE id = ?", b.Id.Local().UnixMilli())
  if err != nil {
    t.Errorf("Failed to query database: %v", err)
    return
  }
  defer rows.Close()

  for rows.Next() {
    var id int64
    var name string
    var price float64
    var priceOne float64
    var quantity float64
    err = rows.Scan(&id, &name, &price, &priceOne, &quantity)
    if err != nil {
      t.Errorf("Failed to scan row: %v", err)
      return
    }
    if id != item.Id.Local().UnixMilli() {
      t.Errorf("Expected ID '%d', got %d", item.Id.Local().UnixMilli(), id)
    }
    if name != item.Name {
      t.Errorf("Expected Name '%s', got %s", item.Name, name)
    }
    if price != item.Price {
      t.Errorf("Expected Price '%f', got %f", item.Price, price)
    }
    if priceOne != item.PriceOne {
      t.Errorf("Expected PriceOne '%f', got %f", item.PriceOne, priceOne)
    }
    if quantity != item.Quantity {
      t.Errorf("Expected Quantity '%f', got %f", item.Quantity, quantity)
    }
  }
}

func TestGetItemsByID(t *testing.T) {
  fmt.Println("Testing GetItemsByID function")
  billRepo, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
  err = billRepo.CreateTables()
  if err != nil {
    t.Errorf("Failed to create tables: %v", err)
    return
  }
  timestamp := time.Now()
  var items []*bill.Item
  item := bill.ItemNew(
    timestamp,
    "item1",
    100.0,
    100.0,
    1.0,
  )
  item2 := bill.ItemNew(
    timestamp,
    "item2",
    102.0,
    102.0,
    2.0,
  )
  items = append(items, item)
  items = append(items, item2)
  err = billRepo.InsertItems(items)
  if err != nil {
    t.Errorf("Failed to insert items: %v", err)
    return
  }

  itemsByID, err := billRepo.GetItemsByID(item.Id.Local().UnixMilli())
  if err != nil {
    t.Errorf("Failed to get items by ID: %v", err)
    return
  }

  if len(itemsByID) != 2 {
    t.Errorf("Expected 2 item, got %d", len(itemsByID))
    return
  }
  if itemsByID[0].Id.Local().UnixMilli() != item.Id.Local().UnixMilli() {
    t.Errorf("Expected ID '%d', got %d", item.Id.Local().UnixMilli(), itemsByID[0].Id.Local().UnixMilli())
  }
  if itemsByID[0].Name != item.Name {
    t.Errorf("Expected Name '%s', got %s", item.Name, itemsByID[0].Name)
  }
  if itemsByID[1].Price != item2.Price {
    t.Errorf("Expected Price '%f', got %f", item2.Price, itemsByID[1].Price)
  }
}

func TestDeleteItems(t *testing.T) {
  fmt.Println("Testing DeleteItems function")
  billRepo, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
  err = billRepo.CreateTables()
  if err != nil {
    t.Errorf("Failed to create tables: %v", err)
    return
  }
  timestamp := time.Now()
  var items []*bill.Item
  item := bill.ItemNew(
    timestamp,
    "item1",
    100.0,
    100.0,
    1.0,
  )
  item2 := bill.ItemNew(
    timestamp,
    "item2",
    102.0,
    102.0,
    2.0,
  )
  items = append(items, item)
  items = append(items, item2)
  err = billRepo.InsertItems(items)
  if err != nil {
    t.Errorf("Failed to insert items: %v", err)
    return
  }

  err = billRepo.DeleteItems(item.Id.Local().UnixMilli())
  if err != nil {
    t.Errorf("Failed to delete items by ID: %v", err)
    return
  }

  itemsFromDb, err := billRepo.GetItemsByID(item.Id.Local().UnixMilli())
  if err != nil {
    t.Errorf("Failed to get items by ID: %v", err)
    return
  }

  if len(itemsFromDb) != 0 {
    t.Errorf("Expected 0 item, got %d", len(itemsFromDb))
    return
  }
}

func TestCheckUniqueItemNames(t *testing.T) {
  fmt.Println("Testing CheckUniqueItemNames function")
  billRepo, err := setUpDB(t)
  if err != nil {
    t.Errorf("Failed to set up database: %v", err)
    return
  }
  err = billRepo.CreateTables()
  if err != nil {
    t.Errorf("Failed to create tables: %v", err)
    return
  }
  timestamp := time.Now()
  var items []*bill.Item
  item := bill.ItemNew(
    timestamp,
    "item1",
    100.0,
    100.0,
    1.0,
  )
  item2 := bill.ItemNew(
    timestamp,
    "item2",
    102.0,
    102.0,
    2.0,
  )

  items = append(items, item)
  items = append(items, item2)
  err = billRepo.InsertItems(items)
  if err != nil {
    t.Errorf("Failed to insert items: %v", err)
    return
  }
  
  err = billRepo.CheckUniqueItemNames()
  if err != nil {
    t.Errorf("Failed to check unique item names: %v", err)
    return
  }

  rows, err := billRepo.DB.Query("SELECT name FROM items_meta")
  if err != nil {
    t.Errorf("Failed to query database: %v", err)
    return
  }
  defer rows.Close()
    
  var names []string
  for rows.Next() {
    var name string
    err = rows.Scan(&name)
    if err != nil {
      t.Errorf("Failed to scan row: %v", err)
      return
    }
    names = append(names, name)
  }
  if item.Name != names[0] {
    t.Errorf("Expected %s, got %s.", item.Name, names[0])
  }
  if item2.Name != names[1] {
    t.Errorf("Expected %s, got %s.", item2.Name, names[1])
  }
}


