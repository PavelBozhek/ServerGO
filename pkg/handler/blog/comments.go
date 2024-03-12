package blog

import (
	"database/sql"
	"net/http"
	"strconv"

	"example.com/server/pkg/models"
	"example.com/server/pkg/repository"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var comment models.Comment
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	err := repository.CreateComment(id, comment.Content, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{Message: "Comment created successfully"})
}

func GetComment(c *gin.Context, db *sql.DB) {

	idStr := c.Param("id")
	var comments []models.Comment

	id, _ := strconv.Atoi(idStr)
	rows, err := repository.GetAllComments(id, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Content); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch comments"})
			return
		}
		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, comments)
}

func DeleteComment(c *gin.Context, db *sql.DB) {
	postIDStr := c.Param("postID")
	commentIDStr := c.Param("commentID")

	commentID, _ := strconv.Atoi(commentIDStr)
	postID, _ := strconv.Atoi(postIDStr)
	err := repository.DeleteComment(commentID, postID, db)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Comment deleted successfully"})
}
