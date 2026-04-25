package repository

import (
	"testing"

	"go-icarros/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCarRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO cars").
		WithArgs(1, "VW", "Gol", 2020, 45000.0).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	repo := &CarRepository{DB: db}
	car := &models.Car{UserID: 1, Marca: "VW", Modelo: "Gol", Ano: 2020, Valor: 45000}

	if err := repo.Create(car); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if car.ID != 10 {
		t.Errorf("esperado ID=10, obtido %d", car.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}

func TestCarRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "user_id", "marca", "modelo", "ano", "valor"}).
		AddRow(1, 1, "VW", "Gol", 2020, 45000.0).
		AddRow(2, 2, "Fiat", "Uno", 2018, 30000.0)

	mock.ExpectQuery("SELECT id, user_id, marca, modelo, ano, valor FROM cars").
		WillReturnRows(rows)

	repo := &CarRepository{DB: db}
	cars, err := repo.FindAll()

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(cars) != 2 {
		t.Errorf("esperado 2 carros, obtido %d", len(cars))
	}
	if cars[0].Marca != "VW" {
		t.Errorf("esperado Marca=VW, obtido %s", cars[0].Marca)
	}
}

func TestCarRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "user_id", "marca", "modelo", "ano", "valor"}).
		AddRow(3, 1, "Honda", "Civic", 2022, 120000.0)

	mock.ExpectQuery("SELECT id, user_id, marca, modelo, ano, valor FROM cars WHERE id").
		WithArgs(3).
		WillReturnRows(rows)

	repo := &CarRepository{DB: db}
	car, err := repo.FindByID(3)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if car.ID != 3 {
		t.Errorf("esperado ID=3, obtido %d", car.ID)
	}
	if car.Marca != "Honda" {
		t.Errorf("esperado Marca=Honda, obtido %s", car.Marca)
	}
}

func TestCarRepository_FindByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "user_id", "marca", "modelo", "ano", "valor"}).
		AddRow(1, 5, "VW", "Gol", 2020, 45000.0).
		AddRow(2, 5, "Fiat", "Palio", 2019, 32000.0)

	mock.ExpectQuery("SELECT id, user_id, marca, modelo, ano, valor FROM cars WHERE user_id").
		WithArgs(5).
		WillReturnRows(rows)

	repo := &CarRepository{DB: db}
	cars, err := repo.FindByUserID(5)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(cars) != 2 {
		t.Errorf("esperado 2 carros, obtido %d", len(cars))
	}
}

func TestCarRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE cars SET").
		WithArgs("Fiat", "Uno", 2021, 28000.0, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &CarRepository{DB: db}
	err = repo.Update(&models.Car{ID: 1, Marca: "Fiat", Modelo: "Uno", Ano: 2021, Valor: 28000})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}

func TestCarRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM cars WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &CarRepository{DB: db}
	if err := repo.Delete(1); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}
