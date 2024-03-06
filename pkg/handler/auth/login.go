package auth

import (
	"database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context, db *sql.DB) {
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
}
