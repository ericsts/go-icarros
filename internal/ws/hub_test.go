package ws

import (
	"testing"
	"time"
)

func TestNewHub_Inicializa(t *testing.T) {
	h := NewHub()

	if h == nil {
		t.Fatal("esperado hub não-nil")
	}
	if h.rooms == nil {
		t.Fatal("rooms deve ser inicializado")
	}
	if h.broadcast == nil {
		t.Fatal("canal broadcast deve ser inicializado")
	}
	if h.register == nil {
		t.Fatal("canal register deve ser inicializado")
	}
	if h.unregister == nil {
		t.Fatal("canal unregister deve ser inicializado")
	}
}

func TestHub_Broadcast_SalaVazia(t *testing.T) {
	h := NewHub()
	go h.Run()
	time.Sleep(5 * time.Millisecond)

	// leilão inexistente — não deve bloquear nem causar pânico
	h.Broadcast(999, map[string]any{"amount": 1000})
	time.Sleep(10 * time.Millisecond)
}

func TestHub_Broadcast_MultiplasSalas(t *testing.T) {
	h := NewHub()
	go h.Run()
	time.Sleep(5 * time.Millisecond)

	// múltiplos broadcasts para salas diferentes sem clientes
	for i := 1; i <= 5; i++ {
		h.Broadcast(i, map[string]any{"auction_id": i, "amount": float64(i * 1000)})
	}
	time.Sleep(10 * time.Millisecond)
}

func TestHub_Broadcast_DadoNaoSerializavel(t *testing.T) {
	h := NewHub()
	go h.Run()
	time.Sleep(5 * time.Millisecond)

	// função não é serializável em JSON — deve logar o erro sem bloquear
	h.Broadcast(1, func() {})
	time.Sleep(10 * time.Millisecond)
}

func TestHub_Broadcast_EstruturaAninhada(t *testing.T) {
	h := NewHub()
	go h.Run()
	time.Sleep(5 * time.Millisecond)

	type bid struct {
		ID        int     `json:"id"`
		AuctionID int     `json:"auction_id"`
		UserID    int     `json:"user_id"`
		Amount    float64 `json:"amount"`
	}
	h.Broadcast(1, bid{ID: 7, AuctionID: 1, UserID: 2, Amount: 46000})
	time.Sleep(10 * time.Millisecond)
}
