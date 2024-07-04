package archive

import (
	"database/sql"
	"testing"
	"time"

	"github.com/segmentio/ksuid"
)

type ChangeBill struct {
	IdOld string
	IdNew string
}

func ChangeBillNew(idOld string, idNew string) *ChangeBill {
	return &ChangeBill{
		IdOld: idOld,
		IdNew: idNew,
	}
}

func TestRunBillMigration(t *testing.T) {
	db, err := sql.Open(
		"sqlite3",
		"/home/qq/Documents/programming/go/billdb/billdbIDchange.db",
	)
	if err != nil {
		t.Errorf("Failed to open sqlite database: %v", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM invoice_tag WHERE invoice_id = '20220726827336';
	DELETE FROM invoice WHERE invoice_id = '20220726827336';`)
	if err != nil {
		t.Errorf("Failed to update ID: %v", err)
		return
	}

	rows, err := db.Query("SELECT invoice_id, invoice_date FROM invoice ORDER BY invoice_date;")
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	layoutDate := "2006-01-02"
	var bills []ChangeBill
	for rows.Next() {
		var (
			idOld  string
			idDate string
		)
		err = rows.Scan(&idOld, &idDate)
		if err != nil {
			t.Errorf("Failed to scan row: %v", err)
			return
		}
		billDateTime, err := time.Parse(layoutDate, idDate)
		if err != nil {
			t.Errorf("Failed to parse date: %v", err)
			return
		}
		idNew, err := ksuid.NewRandomWithTime(billDateTime)
		if err != nil {
			t.Errorf("Failed to generate new ID: %v", err)
			return
		}
		bills = append(bills, *ChangeBillNew(idOld, idNew.String()))
	}
	rows.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to start transaction: %v", err)
		return
	}

	// Disable foreign key constraints using PRAGMA
	_, err = tx.Exec("PRAGMA foreign_keys = OFF;")
	if err != nil {
		t.Error("Error disabling foreign key constraints:", err)
		return
	}

	for _, bill := range bills {
		if err != nil {
			t.Errorf("Failed to generate new ID: %v", err)
			return
		}
		_, err = tx.Exec("UPDATE item SET invoice_id = ? WHERE invoice_id = ?", bill.IdNew, bill.IdOld)
		if err != nil {
			tx.Rollback()
			t.Errorf("Failed to update ID: %v", err)
			return
		}
		_, err = tx.Exec("UPDATE invoice_tag SET invoice_id = ? WHERE invoice_id = ?", bill.IdNew, bill.IdOld)
		if err != nil {
			tx.Rollback()
			t.Errorf("Failed to update ID: %v", err)
			return
		}
		_, err = tx.Exec("UPDATE invoice SET invoice_id = ? WHERE invoice_id = ?", bill.IdNew, bill.IdOld)
		if err != nil {
			tx.Rollback()
			t.Errorf("Failed to update ID: %v", err)
			return
		}
	}

	// Enable foreign key constraints back
	_, err = tx.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		tx.Rollback()
		t.Error("Error enabling foreign key constraints:", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		t.Errorf("Failed to commit transaction: %v", err)
		return
	}
	db.Close()
	t.Errorf("Migration successful")
}

type ItemChange struct {
	IdOld string
	IdNew string
}

func ItemChangeNew(idOld string, idNew string) *ItemChange {
	return &ItemChange{
		IdOld: idOld,
		IdNew: idNew,
	}
}

func TestRunItemMigration(t *testing.T) {
	db, err := sql.Open(
		"sqlite3",
		"/home/qq/Documents/programming/go/billdb/billdbIDchange.db",
	)
	if err != nil {
		t.Errorf("Failed to open sqlite database: %v", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT item_id FROM item;")
	if err != nil {
		t.Errorf("Failed to query database: %v", err)
		return
	}
	defer rows.Close()

	var items []*ItemChange
	for rows.Next() {
		var (
			idOld string
		)
		err = rows.Scan(&idOld)
		if err != nil {
			t.Errorf("Failed to scan row: %v", err)
			return
		}
		idNew := ksuid.New()
		items = append(items, ItemChangeNew(idOld, idNew.String()))
	}
	rows.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Errorf("Failed to start transaction: %v", err)
		return
	}

	// Disable foreign key constraints using PRAGMA
	_, err = tx.Exec("PRAGMA foreign_keys = OFF;")
	if err != nil {
		t.Error("Error disabling foreign key constraints:", err)
		return
	}

	for _, item := range items {
		if err != nil {
			t.Errorf("Failed to generate new ID: %v", err)
			return
		}
		_, err = tx.Exec("UPDATE item SET item_id = ? WHERE item_id = ?", item.IdNew, item.IdOld)
		if err != nil {
			tx.Rollback()
			t.Errorf("Failed to update ID: %v", err)
			return
		}
		_, err = tx.Exec("UPDATE item_photo SET item_id = ? WHERE item_id = ?", item.IdNew, item.IdOld)
		if err != nil {
			tx.Rollback()
			t.Errorf("Failed to update ID: %v", err)
			return
		}
		_, err = tx.Exec("UPDATE item_tag SET item_id = ? WHERE item_id = ?", item.IdNew, item.IdOld)
		if err != nil {
			tx.Rollback()
			t.Errorf("Failed to update ID: %v", err)
			return
		}
	}

	// Enable foreign key constraints back
	_, err = tx.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		tx.Rollback()
		t.Error("Error enabling foreign key constraints:", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		t.Errorf("Failed to commit transaction: %v", err)
		return
	}
	db.Close()
	t.Errorf("Migration successful")
}
