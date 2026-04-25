package service

import (
	"errors"
	"testing"

	"go-icarros/internal/models"
)

type mockCarRepo struct {
	createErr        error
	findAllCars      []models.Car
	findAllErr       error
	findByIDCar      *models.Car
	findByIDErr      error
	findByUserIDCars []models.Car
	findByUserIDErr  error
	updateErr        error
	deleteErr        error
}

func (m *mockCarRepo) Create(_ *models.Car) error          { return m.createErr }
func (m *mockCarRepo) FindAll() ([]models.Car, error)      { return m.findAllCars, m.findAllErr }
func (m *mockCarRepo) FindByID(_ int) (*models.Car, error) { return m.findByIDCar, m.findByIDErr }
func (m *mockCarRepo) FindByUserID(_ int) ([]models.Car, error) {
	return m.findByUserIDCars, m.findByUserIDErr
}
func (m *mockCarRepo) Update(_ *models.Car) error { return m.updateErr }
func (m *mockCarRepo) Delete(_ int) error         { return m.deleteErr }

func TestCarService_Create(t *testing.T) {
	svc := &CarService{Repo: &mockCarRepo{}}

	if err := svc.Create(&models.Car{Marca: "VW"}); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestCarService_Create_PropagaErro(t *testing.T) {
	svc := &CarService{Repo: &mockCarRepo{createErr: errors.New("db error")}}

	if err := svc.Create(&models.Car{}); err == nil {
		t.Fatal("deveria retornar erro")
	}
}

func TestCarService_GetAll(t *testing.T) {
	cars := []models.Car{{ID: 1}, {ID: 2}}
	svc := &CarService{Repo: &mockCarRepo{findAllCars: cars}}

	result, err := svc.GetAll()

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("esperado 2 carros, obtido %d", len(result))
	}
}

func TestCarService_GetByID(t *testing.T) {
	svc := &CarService{Repo: &mockCarRepo{findByIDCar: &models.Car{ID: 3}}}

	car, err := svc.GetByID(3)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if car.ID != 3 {
		t.Errorf("esperado ID=3, obtido %d", car.ID)
	}
}

func TestCarService_GetByID_PropagaErro(t *testing.T) {
	svc := &CarService{Repo: &mockCarRepo{findByIDErr: errors.New("not found")}}

	if _, err := svc.GetByID(99); err == nil {
		t.Fatal("deveria retornar erro")
	}
}

func TestCarService_GetByUserID(t *testing.T) {
	svc := &CarService{Repo: &mockCarRepo{findByUserIDCars: []models.Car{{ID: 1}, {ID: 2}}}}

	cars, err := svc.GetByUserID(1)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(cars) != 2 {
		t.Errorf("esperado 2 carros, obtido %d", len(cars))
	}
}

func TestCarService_Update(t *testing.T) {
	svc := &CarService{Repo: &mockCarRepo{}}

	if err := svc.Update(&models.Car{ID: 1, Marca: "Fiat"}); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestCarService_Delete(t *testing.T) {
	svc := &CarService{Repo: &mockCarRepo{}}

	if err := svc.Delete(1); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
}
