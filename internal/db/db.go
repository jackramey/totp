package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase(db *sql.DB) error {
	// Create tables if they don't exist
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        username TEXT PRIMARY KEY,
        secret TEXT,
        verified BOOLEAN
    );
    `
	if _, err := db.Exec(createTableSQL); err != nil {
		return err
	}

	return nil
}

func InsertUser(db *sql.DB, username, secret string) error {
	insertSQL := "INSERT INTO users (username, secret, verified) VALUES (?, ?, 0)"
	if _, err := db.Exec(insertSQL, username, secret); err != nil {
		return err
	}
	return nil
}

func UpdateVerificationStatus(db *sql.DB, username string) error {
	updateSQL := "UPDATE users SET verified = 1 WHERE username = ?"
	if _, err := db.Exec(updateSQL, username); err != nil {
		return err
	}
	return nil
}

func GetSecretByUsername(db *sql.DB, username string) (string, error) {
	query := "SELECT secret FROM users WHERE username = ?"
	var secret string
	err := db.QueryRow(query, username).Scan(&secret)
	if err != nil {
		return "", err
	}
	return secret, nil
}
