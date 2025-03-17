package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func Connect() *sql.DB {
	db, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	// Setup tables
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetSerializedPassword(db *sql.DB, username string) (string, error) {
	rows, err := db.Query(`SELECT password FROM users WHERE username = ?`, username)
	if err != nil {
		return "", fmt.Errorf("error selecting password for username: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return "", nil
	}

	var password string

	err = rows.Scan(&password)
	if err != nil {
		return "", fmt.Errorf("error retrieving password for username: %w", err)
	}

	if len(password) == 0 {
		return "", fmt.Errorf("password for username is empty")
	}

	return password, nil
}

func AddPassword(username, password string) error {
	db := Connect()
	_, err := db.Exec(`INSERT OR REPLACE INTO users
                       (username, password) VALUES (?, ?)`,
		username,
		password,
	)
	return err
}
