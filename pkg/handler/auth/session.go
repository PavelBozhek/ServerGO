package auth

import (
	"example.com/server/pkg/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Session(c *gin.Context) {

	session := sessions.Default(c)

	userID := session.Get("userID")

	if userID == nil {
		c.JSON(http.StatusOK, models.ErrorResponse{Error: "Session not found"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: fmt.Sprintf("User ID from session: %v", userID)})
}
