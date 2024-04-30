package repository

import (
	bl "billdb/internal/bill"
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

type SqliteBillRepository struct {
	DB *sql.DB
}

func NewSqliteBillRepository(db *sql.DB) *SqliteBillRepository {
	return &SqliteBillRepository{DB: db}
}

// Implementation for getting the database
func (r *SqliteBillRepository) GetDb() *sql.DB {
	return r.DB
}

// Implementation for inserting a bill in the sqlite database
func (r *SqliteBillRepository) InsertBill(bill *bl.Bill) error {
	_, err := r.DB.Exec(`
	INSERT INTO bills (
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
	)
	VALUES (?,?,?,?,?,?,?,?,?,?)
		`,
		bill.GetIdUnix(),
		bill.Name,
		bill.GetDateString(),
		bill.Price,
		bill.GetCurrencyString(),
		bill.ExchangeRate,
		bill.GetCountryString(),
		bill.Tag,
		bill.Link,
		bill.BillText,
	)
	if err != nil {
		log.Error("Error inserting bill into db: ", err)
		return err
	}
	return nil
}

func (r *SqliteBillRepository) GetBills() ([]*bl.Bill, error) {
	query := `
	SELECT
		id, name, dates, price, currency, exchange_rate, country, tag, link, bill
	FROM bills
  LIMIT 20;
  `
	row, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var bills []*bl.Bill
	for row.Next() {
		bill, err := ScanToBill(row)
		if err != nil {
			return nil, err
		}
		bills = append(bills, bill)
	}

	return bills, nil
}

// fetching a bill from the database by ID without items
func (r *SqliteBillRepository) GetBillByID(id int64) (*bl.Bill, error) {
	query := `
	SELECT
		id, name, dates, price, currency, exchange_rate, country, tag, link, bill
	FROM bills
	WHERE id = ?;
  `
	row := r.DB.QueryRow(query, id)
	bill, err := ScanToBill(row)
	if err != nil {
		log.Errorf("Error getting bill by id from db: %v", err)
		return nil, err
	}

	return bill, nil
}

// Implementation for updating a bill in the database
func (r *SqliteBillRepository) UpdateBill(bill *bl.Bill) error {
	_, err := r.DB.Exec(`UPDATE bills SET 
		name = ?,
		dates = ?, 
		price = ?, 
		currency = ?, 
		exchange_rate = ?, 
		country = ?, 
		tag = ?, 
		link = ?, 
		bill = ?
		WHERE id = ?`,
		bill.Name,
		bill.GetDateString(),
		bill.Price,
		bill.GetCurrencyString(),
		bill.ExchangeRate,
		bill.GetCountryString(),
		bill.Tag,
		bill.Link,
		bill.BillText,
		bill.GetIdUnix(),
	)
	if err != nil {
		log.Errorf("Error updating bill by id in db: %v", err)
		return err
	}

	return nil
}

// Implementation for deleting a bill from the database by ID
func (r *SqliteBillRepository) DeleteBill(id int64) error {
	_, err := r.DB.Exec("DELETE FROM bills WHERE id = ?", id)
	if err != nil {
		log.WithField("id", id).Error(
			"Error deleting bill by id from db: ", err)
		return err
	}
	return nil
}

// Implementation for checking unique item names
func (r *SqliteBillRepository) InsertItems(items []*bl.Item) error {
	stmt, err := r.DB.Prepare("INSERT INTO items ( id, name, price, price_one, quantity) VALUES (?,?,?,?,?)")
	if err != nil {
		log.Error("Error inserting bill into db: ", err)
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.Exec(
			item.GetIdUnix(),
			item.Name,
			item.Price,
			item.PriceOne,
			item.Quantity,
		)
		if err != nil {
			log.WithField("itemId", item.GetIdUnix()).Error(
				"Error inserting item into db: ", err)
			return err
		}
	}

	return nil
}

// Implementation for getting an item from the database by ID
func (r *SqliteBillRepository) GetItemsByID(billId int64) ([]*bl.Item, error) {
	rows, err := r.DB.Query(`
	SELECT
		id, 
		name, 
		price, 
		price_one,
		quantity
	FROM items
	WHERE id = ?;
	`, billId)
	if err != nil {
		log.WithField("id", billId).Error(
			"Error getting item by id from db: ", err)
		return nil, err
	}
	defer rows.Close()

	var items []*bl.Item
	for rows.Next() {
		var (
			id       int64
			name     string
			price    float64
			priceOne float64
			quantity float64
		)
		err := rows.Scan(
			&id,
			&name,
			&price,
			&priceOne,
			&quantity,
		)
		if err != nil {
			return nil, err
		}
		item := bl.ItemNew(
			time.UnixMilli(id),
			name,
			price,
			priceOne,
			quantity,
		)
		items = append(items, item)
	}

	return items, nil
}

func (r *SqliteBillRepository) UpdateItems(items []*bl.Item) error {
	return fmt.Errorf("not implemented")
}

// Implementation for updating an item in the database
func (r *SqliteBillRepository) UpdateItem(item *bl.Item) error {
	return fmt.Errorf("not implemented")
}

// Implementation for deleting an item from the database by ID
func (r *SqliteBillRepository) DeleteItems(billId int64) error {
	_, err := r.DB.Exec("DELETE FROM items WHERE id = ?", billId)
	if err != nil {
		log.WithField("id", billId).Error(
			"Error deleting item by id from db: ", err)
		return err
	}
	return nil
}

// createTables creates necessary tables in the SQLite database.
func (r *SqliteBillRepository) CreateTables() error {
	_, err := r.DB.Exec(`
	CREATE TABLE "bills" (
		"id"	INTEGER NOT NULL,
		"name"	TEXT NOT NULL,
		"dates"	TEXT NOT NULL,
		"price"	REAL NOT NULL,
		"tag"	TEXT,
		"currency"	TEXT,
		"exchange_rate"	REAL,
		"country"	TEXT,
		"link"	TEXT,
		"bill"	TEXT,
		PRIMARY KEY("id")
	);

	CREATE TABLE "items" (
		"id"	INTEGER NOT NULL,
		"photo"	BLOB,
		"name"	TEXT,
		"price"	REAL,
		"price_one"	REAL,
		"quantity"	REAL,
		FOREIGN KEY("id") REFERENCES "bills"("id")
	);

	CREATE TABLE "items_meta" (
		"name"	TEXT UNIQUE,
		"tag"	TEXT,
		PRIMARY KEY("name")
	);
  `)
	if err != nil {
		log.Error("Error creating tables in db: ", err)
		return err
	}

	return nil
}

func (r *SqliteBillRepository) CheckDuplicateBill(bill *bl.Bill) ([]*bl.Bill, error) {
	query := `
		SELECT 
      id, name, dates, price, currency, exchange_rate, country, tag, link, bill
		FROM bills
		WHERE
			dates = ? AND price = ? AND currency = ?;
	`
	rows, err := r.DB.Query(query, bill.GetDateString(), bill.Price, bill.GetCurrencyString())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var billArr []*bl.Bill

	for rows.Next() {
		b, err := ScanToBill(rows)
		if err != nil {
			return nil, err
		}
		billArr = append(billArr, b)
	}

	return billArr, nil
}

func (r *SqliteBillRepository) CheckUniqueItemNames() error {
	queryUniqueNames := `
		SELECT DISTINCT i.name
		FROM items as i
		LEFT JOIN items_meta as im ON i.name = im.name
		WHERE im.name IS NULL;
	`
	queryUniqueNamesAdd := "INSERT INTO items_meta (name)\n" + queryUniqueNames

	rows, err := r.DB.Query(queryUniqueNames)
	if err != nil {
		log.Error("Error querying unique item names: ", err)
		return err
	}
	defer rows.Close()

	var names []string

	for rows.Next() {
		var n string
		err := rows.Scan(&n)
		if err != nil {
			log.Error("Error getting item names: ", err)
			return err
		}

		names = append(names, n)
	}

	if len(names) == 0 {
		log.Info("There is no new unique item names.")
		return nil
	}

	_, err = r.DB.Exec(queryUniqueNamesAdd)
	if err != nil {
		log.Error("Error adding new unique item names to the table: ", err)
		return err
	} else {
		log.WithFields(log.Fields{
			"len": len(names),
		}).Info("Successfully added new unique item names to the table.")
		return nil
	}
}

// function to use in place .Scan()
// to scan a row/rows into a bill/bills
// you should use rows.Next()
// and pass the rows to this function
func ScanToBill(row interface{}) (*bl.Bill, error) {
	var (
		Id           int64
		Name         string
		Date         string
		Price        float64
		Currency     string
		ExchangeRate float64
		Country      string
		Tag          string
		Link         string
		BillText     string
	)
	switch r := row.(type) {
	case *sql.Row:
		err := r.Scan(
			&Id,
			&Name,
			&Date,
			&Price,
			&Currency,
			&ExchangeRate,
			&Country,
			&Tag,
			&Link,
			&BillText,
		)
		if err != nil {
			log.Error(
				"Error scaning bill: ", err)
			return nil, err
		}
	case *sql.Rows:
		err := r.Scan(
			&Id,
			&Name,
			&Date,
			&Price,
			&Currency,
			&ExchangeRate,
			&Country,
			&Tag,
			&Link,
			&BillText,
		)
		if err != nil {
			log.Error(
				"Error scaning bill: ", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid type %T", r)
	}
	billDate, err := bl.StringToDate(Date)
	if err != nil {
		log.WithFields(log.Fields{
			"date": Date,
		}).Error(
			"Error parsing bill date from db: ", err)
		return nil, err
	}
	billCurrency, err := bl.StringToCurrency(Currency)
	if err != nil {
		log.WithField("currency", Currency).Error(
			"Error parsing bill currency from db: ", err)
		return nil, err
	}
	billCountry, err := bl.StringToCountry(Country)
	if err != nil {
		log.WithField("country", Country).Error(
			"Error parsing bill country from db: ", err)
		return nil, err
	}

	bill := bl.BillNew(
		bl.UnixToId(Id),
		Name,
		*billDate,
		Price,
		billCurrency,
		ExchangeRate,
		billCountry,
		[]*bl.Item{},
		bl.Tag(Tag),
		Link,
		BillText,
	)

	return bill, nil
}
