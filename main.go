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

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Comment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"postId"`
	Content string `json:"content"`
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
	// Маршруты и обработчики для управления постами и комментариями
	router.POST("/posts", func(c *gin.Context) {
		var post Post
		if err := c.BindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Добавление поста в базу данных
		_, err := db.Exec("INSERT INTO posts (title, content) VALUES ($1, $2)", post.Title, post.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create post"})
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse{Message: "Post created successfully"})
	})

	router.GET("/posts", func(c *gin.Context) {
		var posts []Post
		rows, err := db.Query("SELECT id, title, content FROM posts")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch posts"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch posts"})
				return
			}
			posts = append(posts, post)
		}

		c.JSON(http.StatusOK, posts)
	})

	router.GET("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")
		var post Post
		err := db.QueryRow("SELECT id, title, content FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch post"})
			return
		}

		c.JSON(http.StatusOK, post)
	})

	// Маршрут для обновления поста
	router.PUT("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")
		var post Post
		if err := c.BindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}

		_, err := db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", post.Title, post.Content, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update post"})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{Message: "Post updated successfully"})
	})

	// Маршрут для удаления комментария
	router.DELETE("/posts/:postID/comments/:commentID", func(c *gin.Context) {
		postID := c.Param("postID")
		commentID := c.Param("commentID")

		_, err := db.Exec("DELETE FROM comments WHERE id = $1 AND post_id = $2", commentID, postID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete comment"})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse{Message: "Comment deleted successfully"})
	})

	router.POST("/posts/:id/comments", func(c *gin.Context) {
		id := c.Param("id")
		var comment Comment
		if err := c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
			return
		}

		_, err := db.Exec("INSERT INTO comments (post_id, content) VALUES ($1, $2)", id, comment.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create comment"})
			return
		}

		c.JSON(http.StatusCreated, SuccessResponse{Message: "Comment created successfully"})
	})

	router.GET("/posts/:id/comments", func(c *gin.Context) {
		id := c.Param("id")
		var comments []Comment
		rows, err := db.Query("SELECT id, content FROM comments WHERE post_id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch comments"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var comment Comment
			if err := rows.Scan(&comment.ID, &comment.Content); err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch comments"})
				return
			}
			comments = append(comments, comment)
		}

		c.JSON(http.StatusOK, comments)
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
