package blog

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Comment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"postId"`
	Content string `json:"content"`
}

func CreateComment(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
		return
	}

	_, err := db.Exec("INSERT INTO comments (post_id, content) VALUES ($1, $2)", id, comment.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{Message: "Comment created successfully"})
}

func GetComment(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var comments []Comment
	rows, err := db.Query("SELECT id, content FROM comments WHERE post_id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Content); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch comments"})
			return
		}
		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, comments)
}

func DeleteComment(c *gin.Context, db *sql.DB) {
	postID := c.Param("postID")
	commentID := c.Param("commentID")

	_, err := db.Exec("DELETE FROM comments WHERE id = $1 AND post_id = $2", commentID, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "Comment deleted successfully"})
}
