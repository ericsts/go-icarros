package handler

import (
	"net/http"
	"strconv"
	"strings"

	"go-icarros/internal/service"
	"go-icarros/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	Hub *ws.Hub
}

// ServeAuction faz o upgrade para WebSocket e registra o cliente na sala do leilão.
// Autenticação via query param: /ws/auctions/1?token=<jwt>
func (h *WSHandler) ServeAuction(c *gin.Context) {
	auctionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	tokenStr := c.Query("token")
	if !strings.HasPrefix(tokenStr, "") || tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token obrigatório"})
		return
	}

	claims := &service.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(_ *jwt.Token) (any, error) {
		return []byte("secret_key"), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	h.Hub.RegisterConn(auctionID, conn)
}
