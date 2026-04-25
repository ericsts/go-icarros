package repository

import (
	"database/sql"
	"go-icarros/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.QueryRow(
		"INSERT INTO users(name,email,password,role) VALUES($1,$2,$3,$4) RETURNING id",
		user.Name, user.Email, user.Password, user.Role,
	).Scan(&user.ID)
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.DB.QueryRow(
		"SELECT id, password, role FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Password, &user.Role)

	return &user, err
}
