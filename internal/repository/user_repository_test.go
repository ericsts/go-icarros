package repository

import (
	"testing"

	"go-icarros/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("Eric", "eric@test.com", "hash123", "user").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	repo := &UserRepository{DB: db}
	user := &models.User{Name: "Eric", Email: "eric@test.com", Password: "hash123", Role: "user"}

	if err := repo.Create(user); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("esperado ID=1, obtido %d", user.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "password", "role"}).
		AddRow(1, "hashsenha", "admin")

	mock.ExpectQuery("SELECT id, password, role FROM users WHERE email").
		WithArgs("eric@test.com").
		WillReturnRows(rows)

	repo := &UserRepository{DB: db}
	user, err := repo.FindByEmail("eric@test.com")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("esperado ID=1, obtido %d", user.ID)
	}
	if user.Role != "admin" {
		t.Errorf("esperado role=admin, obtido %s", user.Role)
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(1, "Eric", "eric@test.com", "user").
		AddRow(2, "Ana", "ana@test.com", "admin")

	mock.ExpectQuery("SELECT id, name, email, role FROM users").
		WillReturnRows(rows)

	repo := &UserRepository{DB: db}
	users, err := repo.FindAll()

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("esperado 2 usuários, obtido %d", len(users))
	}
	if users[0].Name != "Eric" {
		t.Errorf("esperado Name=Eric, obtido %s", users[0].Name)
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(5, "Eric", "eric@test.com", "user")

	mock.ExpectQuery("SELECT id, name, email, role FROM users WHERE id").
		WithArgs(5).
		WillReturnRows(rows)

	repo := &UserRepository{DB: db}
	user, err := repo.FindByID(5)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if user.ID != 5 {
		t.Errorf("esperado ID=5, obtido %d", user.ID)
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE users SET").
		WithArgs("Novo Nome", "novo@test.com", "user", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &UserRepository{DB: db}
	err = repo.Update(&models.User{ID: 1, Name: "Novo Nome", Email: "novo@test.com", Role: "user"})

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM users WHERE id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &UserRepository{DB: db}
	if err := repo.Delete(1); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não atendidas: %v", err)
	}
}
