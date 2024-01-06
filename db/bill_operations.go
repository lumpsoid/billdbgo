package db

import (
	"billdb/bill"
	"database/sql"
)

func CheckDuplicateBill(db *sql.DB, bill bill.Bill) error {
	query := `
		SELECT id, name, dates, price, currency, bill
		FROM bills
		WHERE
			dates = ? AND
			price = ? AND
			currency = ?;
	`

	rows, err := db.Query(query, bill.Date, bill.Price, bill.Currency)
	if err != nil {
		return err
	}
	defer rows.Close()

	return err
}
