package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func DatabaseConnection(connString string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}