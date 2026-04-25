package service

import "go-icarros/internal/models"

type CarService struct {
	Repo CarRepo
}

func (s *CarService) Create(car *models.Car) error {
	return s.Repo.Create(car)
}

func (s *CarService) GetAll() ([]models.Car, error) {
	return s.Repo.FindAll()
}

func (s *CarService) GetByID(id int) (*models.Car, error) {
	return s.Repo.FindByID(id)
}

func (s *CarService) GetByUserID(userID int) ([]models.Car, error) {
	return s.Repo.FindByUserID(userID)
}

func (s *CarService) Update(car *models.Car) error {
	return s.Repo.Update(car)
}

func (s *CarService) Delete(id int) error {
	return s.Repo.Delete(id)
}
