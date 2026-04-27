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
		"SELECT id, name, email, password, role FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)

	return &user, err
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		"SELECT id, name, email, role FROM users WHERE id=$1", id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	return &user, err
}

func (r *UserRepository) Update(user *models.User) error {
	_, err := r.DB.Exec(
		"UPDATE users SET name=$1, email=$2, role=$3 WHERE id=$4",
		user.Name, user.Email, user.Role, user.ID,
	)
	return err
}

func (r *UserRepository) UpdatePassword(id int, hashedPassword string) error {
	_, err := r.DB.Exec("UPDATE users SET password=$1 WHERE id=$2", hashedPassword, id)
	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}
