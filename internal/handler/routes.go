package handler

import (
	"database/sql"

	"go-icarros/internal/middleware"
	"go-icarros/internal/repository"
	"go-icarros/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	userRepo := &repository.UserRepository{DB: db}
	userSvc := &service.UserService{Repo: userRepo}
	userH := &UserHandler{Service: userSvc}

	carRepo := &repository.CarRepository{DB: db}
	carSvc := &service.CarService{Repo: carRepo}
	carH := &CarHandler{Service: carSvc}

	// rota pública
	r.POST("/login", userH.Login)

	// rotas de admin (requer auth + role admin)
	admin := r.Group("/")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.POST("/users", userH.Create)
		admin.GET("/users", userH.List)
		admin.GET("/users/:id", userH.GetByID)
		admin.PUT("/users/:id", userH.Update)
		admin.DELETE("/users/:id", userH.Delete)
	}

	// rotas autenticadas (requer auth)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/cars", carH.Create)
		auth.GET("/cars", carH.List)
		auth.GET("/cars/my", carH.GetMyCars)
		auth.GET("/cars/:id", carH.GetByID)
		auth.PUT("/cars/:id", carH.Update)
		auth.DELETE("/cars/:id", carH.Delete)
	}
}
