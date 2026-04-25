package service

import "go-icarros/internal/models"

type UserRepo interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	FindByID(id int) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type CarRepo interface {
	Create(car *models.Car) error
	FindAll() ([]models.Car, error)
	FindByID(id int) (*models.Car, error)
	FindByUserID(userID int) ([]models.Car, error)
	Update(car *models.Car) error
	Delete(id int) error
}
