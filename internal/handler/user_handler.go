package handler

import (
	"database/sql"
	"net/http"

	"go-icarros/internal/models"
	"go-icarros/internal/repository"
	"go-icarros/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	repo := &repository.UserRepository{DB: db}
	svc := &service.UserService{Repo: repo}

	r.POST("/register", func(c *gin.Context) {
		var user models.User
		c.BindJSON(&user)

		err := svc.Register(&user)
		if err != nil {
			c.JSON(500, err)
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.POST("/login", func(c *gin.Context) {
		var input models.User
		c.BindJSON(&input)

		user, err := svc.Login(input.Email, input.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}

		c.JSON(200, user)
	})
}
