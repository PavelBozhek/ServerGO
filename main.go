package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	// Инициализация маршрутизатора Gin
	router := gin.Default()

	// Подключение к базе данных PostgreSQL
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Маршрут для регистрации новых пользователей
	router.POST("/register", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Хэширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.Password = string(hashedPassword)

		// Вставка данных пользователя в таблицу
		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})

	// Запуск сервера на порту 8080
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Функция для создания подключения к базе данных PostgreSQL
func connectDB() (*sql.DB, error) {
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
