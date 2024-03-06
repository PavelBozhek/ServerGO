package main

import (
	"example.com/server/pkg/handler"
	"example.com/server/pkg/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	router := gin.Default()

	store := cookie.NewStore([]byte("your_secret_key_here"))
	router.Use(sessions.Sessions("mysession", store))

	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	handler.InitRoutes(router, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
