package auth

import (
	"example.com/server/pkg/models"
	services "example.com/server/pkg/services/auth_services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Session(c *gin.Context) {

	userID := services.GetSession(c)

	if userID == nil {
		c.JSON(http.StatusOK, models.ErrorResponse{Error: "Session not found"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: fmt.Sprintf("User ID from session: %v", userID)})
}
