package main

import (
	"log"

	"go-icarros/internal/db"
	"go-icarros/internal/models"
	"go-icarros/internal/repository"
	"go-icarros/internal/service"
)

var seeds = []models.User{
	{
		Name:     "Eric Santos",
		Email:    "ericsts@gmail.com",
		Password: "admin123",
		Role:     "admin",
	},
}

func main() {
	database := db.Connect()
	defer database.Close()

	repo := &repository.UserRepository{DB: database}
	svc := &service.UserService{Repo: repo}

	for _, user := range seeds {
		u := user // cópia para evitar captura de ponteiro no loop
		if err := svc.Register(&u); err != nil {
			log.Printf("skip %s: %v\n", u.Email, err)
			continue
		}
		log.Printf("criado: %s (%s)\n", u.Email, u.Role)
	}
}
