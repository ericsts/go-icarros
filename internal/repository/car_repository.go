package repository

import (
	"database/sql"

	"go-icarros/internal/models"
)

type CarRepository struct {
	DB *sql.DB
}

func (r *CarRepository) Create(car *models.Car) error {
	return r.DB.QueryRow(
		"INSERT INTO cars(user_id, marca, modelo, ano, valor) VALUES($1,$2,$3,$4,$5) RETURNING id",
		car.UserID, car.Marca, car.Modelo, car.Ano, car.Valor,
	).Scan(&car.ID)
}

func (r *CarRepository) FindAll() ([]models.Car, error) {
	rows, err := r.DB.Query("SELECT id, user_id, marca, modelo, ano, valor FROM cars")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var c models.Car
		if err := rows.Scan(&c.ID, &c.UserID, &c.Marca, &c.Modelo, &c.Ano, &c.Valor); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}
	return cars, nil
}

func (r *CarRepository) FindByID(id int) (*models.Car, error) {
	var car models.Car
	err := r.DB.QueryRow(
		"SELECT id, user_id, marca, modelo, ano, valor FROM cars WHERE id=$1", id,
	).Scan(&car.ID, &car.UserID, &car.Marca, &car.Modelo, &car.Ano, &car.Valor)
	return &car, err
}

func (r *CarRepository) FindByUserID(userID int) ([]models.Car, error) {
	rows, err := r.DB.Query(
		"SELECT id, user_id, marca, modelo, ano, valor FROM cars WHERE user_id=$1", userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []models.Car
	for rows.Next() {
		var c models.Car
		if err := rows.Scan(&c.ID, &c.UserID, &c.Marca, &c.Modelo, &c.Ano, &c.Valor); err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}
	return cars, nil
}

func (r *CarRepository) Update(car *models.Car) error {
	_, err := r.DB.Exec(
		"UPDATE cars SET marca=$1, modelo=$2, ano=$3, valor=$4 WHERE id=$5",
		car.Marca, car.Modelo, car.Ano, car.Valor, car.ID,
	)
	return err
}

func (r *CarRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM cars WHERE id=$1", id)
	return err
}
