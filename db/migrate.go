package db

import (
	"rest-services/config"
	"log"
)

func Migrate() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	_, err := config.DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	createTokensTable := `
	CREATE TABLE IF NOT EXISTS tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		token TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = config.DB.Exec(createTokensTable)
	if err != nil {
		log.Fatalf("Error creating tokens table: %v", err)
	}

	log.Println("Database migrated successfully.")
}
