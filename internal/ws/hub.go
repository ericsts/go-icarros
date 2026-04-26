package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type client struct {
	auctionID int
	conn      *websocket.Conn
	send      chan []byte
}

type broadcastMsg struct {
	auctionID int
	data      []byte
}

// Hub mantém uma sala por leilão e distribui lances em tempo real.
type Hub struct {
	rooms      map[int]map[*client]bool
	mu         sync.RWMutex
	broadcast  chan broadcastMsg
	register   chan *client
	unregister chan *client
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[int]map[*client]bool),
		broadcast:  make(chan broadcastMsg, 256),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			if h.rooms[c.auctionID] == nil {
				h.rooms[c.auctionID] = make(map[*client]bool)
			}
			h.rooms[c.auctionID][c] = true
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if room := h.rooms[c.auctionID]; room != nil {
				if _, ok := room[c]; ok {
					delete(room, c)
					close(c.send)
				}
				if len(room) == 0 {
					delete(h.rooms, c.auctionID)
				}
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.rooms[msg.auctionID] {
				select {
				case c.send <- msg.data:
				default:
					close(c.send)
					delete(h.rooms[msg.auctionID], c)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast serializa data para JSON e envia para todos os clientes da sala.
func (h *Hub) Broadcast(auctionID int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("ws: erro ao serializar broadcast: %v", err)
		return
	}
	h.broadcast <- broadcastMsg{auctionID: auctionID, data: b}
}

// RegisterConn registra uma nova conexão WebSocket na sala do leilão.
func (h *Hub) RegisterConn(auctionID int, conn *websocket.Conn) {
	c := &client{
		auctionID: auctionID,
		conn:      conn,
		send:      make(chan []byte, 256),
	}
	h.register <- c
	go c.writePump()
	go c.readPump(h)
}

func (c *client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (c *client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()
	for {
		// clientes só recebem — descarta mensagens recebidas e aguarda desconexão
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}
