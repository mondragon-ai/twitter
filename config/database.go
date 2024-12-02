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
// ResetDB recreates all the tables by dropping and creating them again.
func ResetDB(db *sql.DB) error {
	tables := []string{
		`DROP TABLE IF EXISTS mention;`,
		`DROP TABLE IF EXISTS article_url;`,
		`DROP TABLE IF EXISTS tweet_clone;`,
		`DROP TABLE IF EXISTS thread_idea;`,
		`DROP TABLE IF EXISTS tweet_idea;`,
	}

	// Drop all tables
	for _, query := range tables {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error dropping table: %v", err)
		}
	}
	log.Println("All tables dropped successfully.")

	// Create all tables
	creationQueries := []string{
		`CREATE TABLE IF NOT EXISTS mention (
			id SERIAL PRIMARY KEY,
			parent_id TEXT,
			author_id TEXT NOT NULL,
			tweet_id TEXT NOT NULL UNIQUE,
			content TEXT NOT NULL,
			author_name TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS articles (
			id SERIAL PRIMARY KEY,
			url TEXT NOT NULL,
			title TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS clones (
			id SERIAL PRIMARY KEY,
			author_name TEXT,
			tweet TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS threads (
			id SERIAL PRIMARY KEY,
			idea TEXT,
			used_count INT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS ideas (
			id SERIAL PRIMARY KEY,
			idea TEXT,
			used_count INT NOT NULL
		);`,
	}

	for _, query := range creationQueries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}
	log.Println("All tables created successfully.")
	return nil
}

// CreateDB ensures all the tables exist in the database
func CreateDB(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS mention (
			id SERIAL PRIMARY KEY,
			parent_id TEXT,
			author_id TEXT NOT NULL,
			tweet_id TEXT NOT NULL UNIQUE,
			content TEXT NOT NULL,
			author_name TEXT NOT NULL,
			created_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS articles (
			id SERIAL PRIMARY KEY,
			url TEXT NOT NULL,
			title TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS clones (
			id SERIAL PRIMARY KEY,
			author_name TEXT,
			tweet TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS threads (
			id SERIAL PRIMARY KEY,
			idea TEXT,
			used_count INT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS ideas (
			id SERIAL PRIMARY KEY,
			idea TEXT,
			used_count INT NOT NULL
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error ensuring table exists: %v", err)
		}
	}
	log.Println("All tables exist or were created successfully.")
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
