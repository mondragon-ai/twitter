package config

import (
	"database/sql"
	"fmt"
	"log"

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

// RecreateDB deletes the mention table if it exists, then creates a new one
func ResetDB(db *sql.DB) error {
	// Drop the mention table if it exists
	_, err := db.Exec(`DROP TABLE IF EXISTS mention;`)
	if err != nil {
		return fmt.Errorf("error deleting mention table: %v", err)
	}
	log.Println("mention table deleted successfully.")

	// Create a new mention table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS mention (
			id SERIAL PRIMARY KEY,
			parent_id TEXT,               	-- Conversation ID (Parent Tweet ID)
			author_id TEXT NOT NULL,      	-- Author ID
			tweet_id TEXT NOT NULL UNIQUE, 	-- Tweet ID (Unique identifier for the tweet)
			content TEXT NOT NULL,        	-- Tweet content
			author_name TEXT NOT NULL,    	-- Author's name
			created_at TEXT NOT NULL 		-- Timestamp for when the mention was recorded
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating mention table: %v", err)
	}
	log.Println("mention table created successfully.")
	return nil
}

// CreateDB ensures the mention table exists in the database
func CreateDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS mention (
			id SERIAL PRIMARY KEY,
			parent_id TEXT,               	-- Conversation ID (Parent Tweet ID)
			author_id TEXT NOT NULL,      	-- Author ID
			tweet_id TEXT NOT NULL UNIQUE, 	-- Tweet ID (Unique identifier for the tweet)
			content TEXT NOT NULL,        	-- Tweet content
			author_name TEXT NOT NULL,    	-- Author's name
			created_at TEXT NOT NULL 		-- Timestamp for when the mention was recorded
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating mention table: %v", err)
	}

	fmt.Println("Mention table exists or was created successfully.")
	return nil
}

// func CreateDB(db *sql.DB) error {
// 	tables := []string{
// 		`CREATE TABLE IF NOT EXISTS mention (
// 			id SERIAL PRIMARY KEY,
// 			content TEXT NOT NULL,
// 			author TEXT NOT NULL,
// 			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 		);`,
// 		`CREATE TABLE IF NOT EXISTS another_table (
// 			id SERIAL PRIMARY KEY,
// 			name TEXT NOT NULL
// 		);`,
// 	}

// 	for _, table := range tables {
// 		_, err := db.Exec(table)
// 		if err != nil {
// 			return fmt.Errorf("error creating table: %v", err)
// 		}
// 	}

// 	log.Println("All tables exist or were created successfully.")
// 	return nil
// }
