package services

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func MakeSession(c *gin.Context, userID int) {
	session := sessions.Default(c)
	session.Set("userID", userID)
	session.Save()
}

func GetSession(c *gin.Context) interface{} {
	session := sessions.Default(c)
	userID := session.Get("userID")
	return userID
}
