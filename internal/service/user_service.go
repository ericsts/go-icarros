package service

import (
	"go-icarros/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo UserRepo
}

func (s *UserService) Register(user *models.User) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hashed)

	return s.Repo.Create(user)
}

func (s *UserService) Login(email, password string) (*models.User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.Repo.FindAll()
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	return s.Repo.FindByID(id)
}

func (s *UserService) Update(user *models.User) error {
	return s.Repo.Update(user)
}

func (s *UserService) UpdateProfile(user *models.User, newPassword string) error {
	current, err := s.Repo.FindByID(user.ID)
	if err != nil {
		return err
	}
	user.Role = current.Role
	if err := s.Repo.Update(user); err != nil {
		return err
	}
	if newPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
		if err != nil {
			return err
		}
		return s.Repo.UpdatePassword(user.ID, string(hashed))
	}
	return nil
}

func (s *UserService) Delete(id int) error {
	return s.Repo.Delete(id)
}
