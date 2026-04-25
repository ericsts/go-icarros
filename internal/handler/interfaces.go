package handler

import "go-icarros/internal/models"

type UserSvc interface {
	Register(user *models.User) error
	Login(email, password string) (*models.User, error)
	GetAll() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type CarSvc interface {
	Create(car *models.Car) error
	GetAll() ([]models.Car, error)
	GetByID(id int) (*models.Car, error)
	GetByUserID(userID int) ([]models.Car, error)
	Update(car *models.Car) error
	Delete(id int) error
}
