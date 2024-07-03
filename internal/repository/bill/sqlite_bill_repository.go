package repository

import (
	bl "billdb/internal/bill"
	"billdb/internal/bill/country"
	"billdb/internal/bill/currency"
	"billdb/internal/bill/item"
	"billdb/internal/bill/tag"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
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
	// start transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO invoice (
			invoice_id, 
			invoice_name, 
			invoice_date, 
			invoice_price, 
			invoice_currency, 
			invoice_country, 
			invoice_link, 
			invoice_text
		)
		VALUES (?,?,?,?,?,?,?,?)`,
		bill.Id,
		bill.Name,
		bill.GetDateString(),
		bill.Price,
		bill.GetCurrencyString(),
		// TODO exchange rate system
		// bill.ExchangeRate,
		bill.GetCountryString(),
		bill.Link,
		bill.BillText,
	)
	if err != nil {
		return err
	}
	if len(bill.Tag) > 0 {
		var tagID int64
		err = tx.QueryRow("SELECT tag_id FROM tag WHERE tag_name = ?", bill.Tag).Scan(&tagID)
		if err != nil && err != sql.ErrNoRows {
			tx.Rollback()
			return err
		}

		if tagID == 0 { // Tag does not exist, insert it
			result, err := tx.Exec("INSERT INTO tag (tag_name) VALUES (?)", bill.Tag)
			if err != nil {
				tx.Rollback()
				return err
			}
			tagID, err = result.LastInsertId()
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		// Link invoice and tag
		_, err = tx.Exec("INSERT INTO invoice_tag (invoice_id, tag_id) VALUES (?, ?)", bill.Id, tagID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// Commit transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// fetching a bill from the database by ID without items
func (r *SqliteBillRepository) GetBillByID(id string) (*bl.Bill, error) {
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
	row := r.DB.QueryRow(query, id)
	bill, err := ScanToBill(row)
	if err != nil {
		return nil, err
	}

	return bill, nil
}

// Implementation for updating a bill in the database
func (r *SqliteBillRepository) UpdateBill(bill *bl.Bill) error {
	// start transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	result, err := tx.Exec(`UPDATE invoice
		SET 
			invoice_name = ?,
			invoice_date = ?, 
			invoice_price = ?, 
			invoice_currency = ?, 
			invoice_country = ?, 
			invoice_link = ?, 
			invoice_text = ?
		WHERE invoice_id = ?`,
		bill.Name,
		bill.GetDateString(),
		bill.Price,
		bill.GetCurrencyString(),
		// TODO exchange rate system
		// bill.ExchangeRate,
		bill.GetCountryString(),
		bill.Link,
		bill.BillText,
		bill.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	// check if provided ID was in the db
	// and was there any change after our UPDATE
	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsUpdated == 0 {
		tx.Rollback()
		return fmt.Errorf("no rows affected")
	}

	// Prepare the INSERT statement to insert tag_name
	// if it doesn't already exist
	insertQuery := `
	INSERT INTO tag (tag_name)
	SELECT ? 
	WHERE NOT EXISTS (SELECT 1 FROM tag WHERE tag_name = ?);
	`

	// Execute the INSERT statement
	_, err = tx.Exec(insertQuery, bill.Tag, bill.Tag)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting tag: %w", err)
	}

	// Now retrieve the tag_id using the SELECT statement
	selectQuery := `
	SELECT tag_id FROM tag WHERE tag_name = ?;
	`

	var tagID int64
	err = tx.QueryRow(selectQuery, bill.Tag).Scan(&tagID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting tag_id: %w", err)
	}

	// Link invoice and tag
	_, err = tx.Exec(
		`INSERT INTO invoice_tag (invoice_id, tag_id) VALUES (?, ?)
    ON CONFLICT(invoice_id) DO UPDATE SET tag_id = excluded.tag_id`,
		bill.Id,
		tagID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Implementation for deleting a bill from the database by ID
func (r *SqliteBillRepository) DeleteBill(id string) error {
	_, err := r.DB.Exec(
		`DELETE FROM invoice WHERE invoice_id = ?;
		DELETE FROM invoice_tag WHERE invoice_id = ?;`,
		id,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

// Implementation for checking unique item names
func (r *SqliteBillRepository) InsertItems(items []*item.Item) error {
	stmt, err := r.DB.Prepare(
		"INSERT INTO item ( item_id, invoice_id, item_name, item_price, item_price_one, item_quantity) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.Exec(
			item.ItemId,
			item.BillId,
			item.Name,
			item.Price,
			item.PriceOne,
			item.Quantity,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Implementation for getting an item from the database by ID
func (r *SqliteBillRepository) GetItemsByID(billId string) ([]*item.Item, error) {
	rows, err := r.DB.Query(`SELECT
			item_id, 
			invoice_id, 
			item_name, 
			item_price, 
			item_price_one,
			item_quantity
		FROM item
		WHERE invoice_id = ?;`,
		billId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*item.Item
	for rows.Next() {
		var (
			itemId   string
			billId   string
			name     string
			price    float64
			priceOne float64
			quantity float64
		)
		err := rows.Scan(
			&itemId,
			&billId,
			&name,
			&price,
			&priceOne,
			&quantity,
		)
		if err != nil {
			return nil, err
		}
		it := item.New(
			itemId,
			billId,
			name,
			price,
			priceOne,
			quantity,
		)
		items = append(items, it)
	}

	return items, nil
}

func (r *SqliteBillRepository) UpdateItems(items []*item.Item) error {
	return fmt.Errorf("not implemented")
}

// Implementation for updating an item in the database
func (r *SqliteBillRepository) UpdateItem(item *item.Item) error {
	return fmt.Errorf("not implemented")
}

// Implementation for deleting an item from the database by ID
func (r *SqliteBillRepository) DeleteItems(items []*item.Item) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	for _, item := range items {
		_, err := tx.Exec(
			`DELETE FROM item WHERE item_id = ?;
			DELETE FROM item_tag WHERE item_id = ?;`,
			item.ItemId,
			item.ItemId,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *SqliteBillRepository) GetCountries() ([]string, error) {
	rows, err := r.DB.Query("SELECT DISTINCT invoice_country FROM invoice;")
	if err != nil {
		log.Error("Error getting countries from db: ", err)
		return nil, err
	}
	defer rows.Close()

	countries := []string{}
	for rows.Next() {
		var country string
		err := rows.Scan(&country)
		if err != nil {
			log.Error("Error scanning countries from db: ", err)
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func (r *SqliteBillRepository) GetCurrencies() ([]string, error) {
	rows, err := r.DB.Query("SELECT DISTINCT invoice_currency FROM invoice;")
	if err != nil {
		log.Error("Error getting currencies from db: ", err)
		return nil, err
	}
	defer rows.Close()

	currencies := []string{}
	for rows.Next() {
		var currency string
		err := rows.Scan(&currency)
		if err != nil {
			log.Error("Error scanning currencies from db: ", err)
			return nil, err
		}
		currencies = append(currencies, currency)
	}
	return currencies, nil
}

func (r *SqliteBillRepository) GetTags() ([]string, error) {
	rows, err := r.DB.Query(`SELECT tag_name FROM tag;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// createTables creates necessary tables in the SQLite database.
func (r *SqliteBillRepository) ApplyMigration(sqlFilePath string) error {

	// Read migration sql file
	sqlFile, err := os.ReadFile(sqlFilePath)
	if err != nil {
		log.Fatalf("Error reading SQL file: %v\n", err)
	}

	// Split SQL file content into individual statements based on ';'
	sqlStatements := strings.Split(string(sqlFile), ";")

	// Iterate over each SQL statement and execute it
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err := r.DB.Exec(stmt)
		if err != nil {
			log.Fatalf("Error executing SQL statement: %v\nSQL: %s\n", err, stmt)
		}
	}

	_, err = r.DB.Exec(`INSERT INTO migration (name) VALUES (?)`, filepath.Base(sqlFilePath))
	if err != nil {
		return err
	}

	return nil
}

func (r *SqliteBillRepository) CheckDuplicateBill(bill *bl.Bill) (int, error) {
	query := `SELECT invoice_id
		FROM invoice
		WHERE invoice_date = ?
			AND invoice_price = ?
			AND invoice_currency = ?;`
	rows, err := r.DB.Query(
		query,
		bill.GetDateString(),
		bill.Price,
		bill.GetCurrencyString(),
	)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var billArr []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return -1, err
		}
		billArr = append(billArr, id)
	}

	return len(billArr), nil
}

// function to use in place .Scan()
//
// to scan a row/rows into a bill/bills
//
// you should use rows.Next()
// and pass the rows to this function
func ScanToBill(row interface{}) (*bl.Bill, error) {
	var (
		Id       string
		Name     string
		Date     string
		Price    float64
		Currency string
		Country  string
		Tag      *string
		Link     string
	)
	switch r := row.(type) {
	case *sql.Row:
		err := r.Scan(
			&Id,
			&Name,
			&Date,
			&Price,
			&Currency,
			&Country,
			&Tag,
			&Link,
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
			&Country,
			&Tag,
			&Link,
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
		return nil, err
	}
	billCurrency, err := currency.Parse(Currency)
	if err != nil {
		return nil, err
	}
	billCountry, err := country.Parse(Country)
	if err != nil {
		return nil, err
	}

	bill := bl.New(
		Id,
		Name,
		*billDate,
		Price,
		billCurrency,
		// TODO exchange rate system
		// ExchangeRate,
		billCountry,
		[]*item.Item{},
		tag.New(Tag),
		Link,
		"",
	)

	return bill, nil
}
