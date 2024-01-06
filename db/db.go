package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// InitializeDB initializes the SQLite database.
func InitializeDB(dbPath string) (*sql.DB, error) {
	// Open an SQLite database connection
	db, err := OpenDb(dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Create necessary tables and perform other initialization steps if needed
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// createTables creates necessary tables in the SQLite database.
func createTables(db *sql.DB) error {
	// Placeholder implementation, replace with actual table creation SQL statements
	// This is just an example; adjust according to your schema
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY,
            username TEXT NOT NULL,
            password TEXT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY,
            user_id INTEGER,
            amount INTEGER,
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
    `)

	return err
}

func OpenDb(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CloseDB closes the database connection.
func CloseDB(db *sql.DB) error {
	// Placeholder implementation, replace with actual database closing code
	return nil
}
