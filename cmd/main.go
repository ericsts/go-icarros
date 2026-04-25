package main

import (
	"log"

	"go-icarros/internal/db"
	"go-icarros/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	database := db.Connect()

	r := gin.Default()

	handler.RegisterRoutes(r, database)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
