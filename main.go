package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	// Инициализация маршрутизатора Gin
	router := gin.Default()

	store := cookie.NewStore([]byte("your_secret_key_here"))
	router.Use(sessions.Sessions("mysession", store))

	// Подключение к базе данных PostgreSQL
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Маршрут для регистрации новых пользователей
	router.POST("/register", func(c *gin.Context) {
		var req RegistrationRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Хэширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
			return
		}

		user := User{Username: req.Username, Email: req.Email, Password: string(hashedPassword)}

		// Вставка данных пользователя в таблицу
		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to register user"})
			return
		}
		// Получаем сеанс
		session := sessions.Default(c)
		// Устанавливаем пользовательский идентификатор в сеанс
		session.Set("userID", user.ID)
		// Сохраняем сеанс
		session.Save()

		c.JSON(http.StatusOK, SuccessResponse{Message: "User registered successfully"})
	})

	router.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}
		var user User
		err := db.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1", req.Username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
			return
		}
		// Получаем сеанс
		session := sessions.Default(c)
		// Устанавливаем пользовательский идентификатор в сеанс
		session.Set("userID", user.ID)
		// Сохраняем сеанс
		session.Save()

		c.JSON(http.StatusOK, SuccessResponse{Message: "Login successful"})
	})

	// Маршрут для получения информации о сессии пользователя
	router.GET("/session-info", func(c *gin.Context) {
		// Получаем сеанс
		session := sessions.Default(c)
		// Получаем значение из сессии (например, идентификатор пользователя)
		userID := session.Get("userID")

		// Проверяем, есть ли значение в сессии
		if userID == nil {
			c.JSON(http.StatusOK, ErrorResponse{Error: "Session not found"})
			return
		}

		// Если значение найдено, отображаем его
		c.JSON(http.StatusOK, SuccessResponse{Message: fmt.Sprintf("User ID from session: %v", userID)})
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
