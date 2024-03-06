package auth

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Session(c *gin.Context) {
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
}
