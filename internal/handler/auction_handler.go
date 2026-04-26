package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuctionHandler struct {
	Service AuctionSvc
}

func (h *AuctionHandler) List(c *gin.Context) {
	auctions, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, auctions)
}

func (h *AuctionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	auction, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "leilão não encontrado"})
		return
	}
	c.JSON(http.StatusOK, auction)
}

func (h *AuctionHandler) PlaceBid(c *gin.Context) {
	auctionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	userID := c.GetInt("user_id")

	var input struct {
		Amount float64 `json:"amount"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bid, err := h.Service.PlaceBid(auctionID, userID, input.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, bid)
}

func (h *AuctionHandler) GetBids(c *gin.Context) {
	auctionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	bids, err := h.Service.GetBids(auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bids)
}
