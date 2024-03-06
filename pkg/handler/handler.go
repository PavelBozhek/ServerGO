package handler

import (
	"database/sql"
	"example.com/server/pkg/handler/auth"
	"example.com/server/pkg/handler/blog"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine, db *sql.DB) {
	router.POST("/register", func(c *gin.Context) {
		auth.Register(c, db)
	})
	router.POST("/posts", func(c *gin.Context) {
		blog.CreatePost(c, db)
	})
	router.GET("/posts", func(c *gin.Context) {
		blog.GetPost(c, db)
	})
	router.GET("/posts/:id", func(c *gin.Context) {
		blog.GetPostById(c, db)
	})
	router.PUT("/posts/:id", func(c *gin.Context) {
		blog.UpdatePost(c, db)
	})
	router.DELETE("/posts/:postID/comments/:commentID", func(c *gin.Context) {
		blog.DeletePost(c, db)
	})
	router.POST("/posts/:id/comments", func(c *gin.Context) {
		blog.CreateComment(c, db)
	})
	router.GET("/posts/:id/comments", func(c *gin.Context) {
		blog.GetComment(c, db)
	})
	router.POST("/login", func(c *gin.Context) {
		auth.Login(c, db)
	})
	router.GET("/session-info", func(c *gin.Context) {
		auth.Session(c)
	})
}
