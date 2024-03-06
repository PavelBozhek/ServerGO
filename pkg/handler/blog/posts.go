package blog

import (
	"database/sql"
	"example.com/server/pkg/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreatePost(c *gin.Context, db *sql.DB) {
	var post models.Post
	if err := c.BindJSON(&post); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	fmt.Println("Received post:", post)

	_, err := db.Exec("INSERT INTO posts (title, content, owner) VALUES ($1, $2, $3)", post.Title, post.Content, post.Owner)
	if err != nil {
		fmt.Println("Error creating post:", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create post"})
		return
	}

	fmt.Println("Post created successfully:", post)
	c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Post created successfully"})
}

func GetPost(c *gin.Context, db *sql.DB) {
	var posts []models.Post
	rows, err := db.Query("SELECT id, title, content FROM posts")
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
	id := c.Param("id")
	var post models.Post
	err := db.QueryRow("SELECT id, title, content FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func UpdatePost(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var post models.Post
	if err := c.BindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	_, err := db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", post.Title, post.Content, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Post updated successfully"})
}
func DeletePost(c *gin.Context, db *sql.DB) {
	postID := c.Param("postID")
	commentID := c.Param("commentID")

	_, err := db.Exec("DELETE FROM comments WHERE id = $1 AND post_id = $2", commentID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Comment deleted successfully"})
}
