package auth

import (
	"database/sql"
	"net/http"

	"example.com/server/pkg/models"
	"example.com/server/pkg/repository"
	services "example.com/server/pkg/services/auth_services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context, db *sql.DB) {
	var req models.LoginRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	user, err := repository.LoginUser(req.Username, db)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	services.MakeSession(c, user.ID)

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Login successful"})
}
