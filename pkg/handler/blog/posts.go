package blog

import (
	"database/sql"
	services "example.com/server/pkg/services/auth_services"
	"fmt"
	"net/http"
	"strconv"

	"example.com/server/pkg/models"
	"example.com/server/pkg/repository"
	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context, db *sql.DB) {
	var post models.Post
	if err := c.BindJSON(&post); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	err := repository.CreatePost(post.Title, post.Content, services.GetSession(c).(int), db)

	if err != nil {
		fmt.Println("Error creating post:", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Post created successfully"})
}

func GetPost(c *gin.Context, db *sql.DB) {
	var posts []models.Post
	rows, err := repository.GetAllPost(db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch posts"})
			return
		}
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

func GetPostById(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	post, err := repository.GetPostByID(id, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func UpdatePost(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	var post models.Post
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	err := repository.UpdatePost(post, db, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Post updated successfully"})
}
func DeletePost(c *gin.Context, db *sql.DB) {
	postIDStr := c.Param("postID")
	postID, _ := strconv.Atoi(postIDStr)
	err := repository.DeletePost(postID, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Comment deleted successfully"})
}
