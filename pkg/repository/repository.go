package repository

import (
	"database/sql"
	"fmt"
)

// Функция для создания подключения к базе данных PostgreSQL
func ConnectDB() (*sql.DB, error) {
	connectionInfo := "host=localhost port=5432 user=postgres password=1234 dbname=mydatabase sslmode=disable"

	db, err := sql.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	// Проверка подключения к базе данных
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
func CreateUser(username, email, password string, db *sql.DB) (int, error) {
	var UserID int
	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", username, email, password)
	if err != nil {
		return 0, err
	}
	return UserID, nil
}
