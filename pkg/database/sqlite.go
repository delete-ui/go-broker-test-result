package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	if err := createTables(db); err != nil {
		log.Printf("Error creating tables: %v", err)
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queryTrades := `
		CREATE TABLE IF NOT EXISTS trades_q (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account TEXT NOT NULL,
			symbol TEXT NOT NULL,
			volume REAL NOT NULL,
			open REAL NOT NULL,
			close REAL NOT NULL,
			side TEXT NOT NULL,
			processed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CHECK (volume > 0),
			CHECK (open > 0),
			CHECK (close > 0),
			CHECK (side IN ('buy', 'sell'))
		);
	`

	queryStats := `
		CREATE TABLE IF NOT EXISTS account_stats (
			account TEXT PRIMARY KEY,
			trades INTEGER NOT NULL DEFAULT 0,
			profit REAL NOT NULL DEFAULT 0
		);
	`

	if _, err := db.Exec(queryTrades); err != nil {
		log.Printf("Error creating tables: %v", err)
		return err
	}

	if _, err := db.Exec(queryStats); err != nil {
		log.Printf("Error creating stats table: %v", err)
		return err
	}

	return nil

}
