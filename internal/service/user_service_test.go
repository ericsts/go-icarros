package service

import (
	"errors"
	"testing"

	"go-icarros/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	createErr          error
	findByEmailUser    *models.User
	findByEmailErr     error
	findAllUsers       []models.User
	findAllErr         error
	findByIDUser       *models.User
	findByIDErr        error
	updateErr          error
	updatePasswordErr  error
	deleteErr          error
}

func (m *mockUserRepo) Create(_ *models.User) error { return m.createErr }
func (m *mockUserRepo) FindByEmail(_ string) (*models.User, error) {
	return m.findByEmailUser, m.findByEmailErr
}
func (m *mockUserRepo) FindAll() ([]models.User, error)             { return m.findAllUsers, m.findAllErr }
func (m *mockUserRepo) FindByID(_ int) (*models.User, error)        { return m.findByIDUser, m.findByIDErr }
func (m *mockUserRepo) Update(_ *models.User) error                 { return m.updateErr }
func (m *mockUserRepo) UpdatePassword(_ int, _ string) error        { return m.updatePasswordErr }
func (m *mockUserRepo) Delete(_ int) error                          { return m.deleteErr }

func TestUserService_Register_HashaSenha(t *testing.T) {
	svc := &UserService{Repo: &mockUserRepo{}}
	user := &models.User{Password: "plaintext"}

	err := svc.Register(user)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if user.Password == "plaintext" {
		t.Error("a senha deveria ter sido hasheada")
	}
}

func TestUserService_Register_PropagaErroDoRepo(t *testing.T) {
	svc := &UserService{Repo: &mockUserRepo{createErr: errors.New("db error")}}

	err := svc.Register(&models.User{Password: "abc"})

	if err == nil {
		t.Fatal("deveria retornar erro")
	}
}

func TestUserService_Login_Sucesso(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.MinCost)
	svc := &UserService{
		Repo: &mockUserRepo{
			findByEmailUser: &models.User{ID: 1, Role: "user", Password: string(hash)},
		},
	}

	user, err := svc.Login("test@example.com", "senha123")

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("esperado ID=1, obtido %d", user.ID)
	}
}

func TestUserService_Login_SenhaErrada(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("certa"), bcrypt.MinCost)
	svc := &UserService{
		Repo: &mockUserRepo{
			findByEmailUser: &models.User{Password: string(hash)},
		},
	}

	_, err := svc.Login("test@example.com", "errada")

	if err == nil {
		t.Fatal("deveria retornar erro para senha errada")
	}
}

func TestUserService_Login_UsuarioNaoEncontrado(t *testing.T) {
	svc := &UserService{
		Repo: &mockUserRepo{findByEmailErr: errors.New("not found")},
	}

	_, err := svc.Login("naoexiste@example.com", "qualquer")

	if err == nil {
		t.Fatal("deveria retornar erro quando usuário não existe")
	}
}

func TestUserService_GetAll(t *testing.T) {
	users := []models.User{{ID: 1}, {ID: 2}}
	svc := &UserService{Repo: &mockUserRepo{findAllUsers: users}}

	result, err := svc.GetAll()

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("esperado 2 usuários, obtido %d", len(result))
	}
}

func TestUserService_GetByID(t *testing.T) {
	svc := &UserService{Repo: &mockUserRepo{findByIDUser: &models.User{ID: 5}}}

	user, err := svc.GetByID(5)

	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if user.ID != 5 {
		t.Errorf("esperado ID=5, obtido %d", user.ID)
	}
}

func TestUserService_Update(t *testing.T) {
	svc := &UserService{Repo: &mockUserRepo{}}

	if err := svc.Update(&models.User{ID: 1, Name: "Atualizado"}); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestUserService_Update_PropagaErro(t *testing.T) {
	svc := &UserService{Repo: &mockUserRepo{updateErr: errors.New("db error")}}

	if err := svc.Update(&models.User{ID: 1}); err == nil {
		t.Fatal("deveria retornar erro")
	}
}

func TestUserService_Delete(t *testing.T) {
	svc := &UserService{Repo: &mockUserRepo{}}

	if err := svc.Delete(1); err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func TestUserService_Delete_PropagaErro(t *testing.T) {
	svc := &UserService{Repo: &mockUserRepo{deleteErr: errors.New("db error")}}

	if err := svc.Delete(1); err == nil {
		t.Fatal("deveria retornar erro")
	}
}
