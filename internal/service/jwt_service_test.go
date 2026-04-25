package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken_RetornaToken(t *testing.T) {
	token, err := GenerateToken(1, "admin")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if token == "" {
		t.Fatal("token não pode ser vazio")
	}
}

func TestGenerateToken_ClaimsCorretos(t *testing.T) {
	tokenStr, err := GenerateToken(42, "user")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(_ *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		t.Fatalf("token deveria ser válido, erro: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("esperado UserID=42, obtido %d", claims.UserID)
	}
	if claims.Role != "user" {
		t.Errorf("esperado Role=user, obtido %s", claims.Role)
	}
	if !claims.ExpiresAt.Time.After(time.Now()) {
		t.Error("token não deve estar expirado")
	}
}

func TestGenerateToken_TokensDiferentes(t *testing.T) {
	t1, _ := GenerateToken(1, "user")
	t2, _ := GenerateToken(2, "admin")
	if t1 == t2 {
		t.Error("tokens de usuários diferentes devem ser distintos")
	}
}
