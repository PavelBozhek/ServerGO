package auth

import (
	"database/sql"
	"example.com/server/pkg/models"
	"example.com/server/pkg/repository"
	services "example.com/server/pkg/services/auth_services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Register(c *gin.Context, db *sql.DB) {
	var req models.RegistrationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}
	// to services
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to hash password"})
		return
	}
	// to services
	userID, err := repository.CreateUser(req.Username, req.Email, string(hashedPassword), db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to register user"})
		return
	}

	services.MakeSession(c, userID)

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "User registered successfully"})
}
