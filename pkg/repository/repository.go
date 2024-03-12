package repository

import (
	"database/sql"
	"example.com/server/pkg/models"
	"fmt"
	"time"
)

func ConnectDB() (*sql.DB, error) {
	connectionInfo := "host=localhost port=5432 user=postgres password=1234 dbname=mydatabase sslmode=disable"

	db, err := sql.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func CreateUser(username, email, password string, db *sql.DB) (int, error) {
	var UserID int
	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", username, email, password)
	if err != nil {
		return 0, err
	}
	return UserID, nil
}

func LoginUser(userName string, db *sql.DB) (models.User, error) {
	var user models.User

	err := db.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1", userName).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	return user, err
}

func CreateComment(id, content string, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO comments (post_id, content, time) VALUES ($1, $2, $3)", id, content, time.Now())
	return err
}

func GetAllComments(id int, db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT id, content FROM comments WHERE post_id = $1", id)
	return rows, err
}

func DeleteComment(commentID, postID int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM comments WHERE id = $1 AND post_id = $2", commentID, postID)
	return err
}

func CreatePost(title string, content string, owner int, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO posts (title, content, owner, time) VALUES ($1, $2, $3, $4)", title, content, owner, time.Now())
	return err
}

func GetAllPost(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT id, title, content FROM posts")
	return rows, err
}

func GetPostByID(id int, db *sql.DB) (models.Post, error) {
	var post models.Post
	err := db.QueryRow("SELECT id, title, content FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content)
	return post, err
}

func UpdatePost(post models.Post, db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", post.Title, post.Content, id)
	return err
}

func DeletePost(postID int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM posts WHERE id = $1", postID)
	return err
}
