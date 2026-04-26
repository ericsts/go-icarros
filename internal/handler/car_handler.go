package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go-icarros/internal/models"

	"github.com/gin-gonic/gin"
)

type CarHandler struct {
	Service    CarSvc
	AuctionSvc AuctionSvc
	Logger     Logger
	Publisher  Publisher
}

// createCarInput agrega os campos do carro e do leilão num único body.
type createCarInput struct {
	models.Car
	AuctionEndsAt time.Time `json:"auction_ends_at"`
	MinBid        float64   `json:"min_bid"`
}

func (h *CarHandler) Create(c *gin.Context) {
	userID := c.GetInt("user_id")

	var input createCarInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Car.UserID = userID
	if err := h.Service.Create(&input.Car); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.AuctionSvc.CreateForCar(input.Car.ID, input.AuctionEndsAt, input.MinBid); err != nil {
		h.Logger.Warn("auction.create_failed", "falha ao criar leilão para o carro", map[string]any{
			"car_id": input.Car.ID,
			"error":  err.Error(),
		})
	}

	meta, _ := json.Marshal(map[string]any{
		"car_id": input.Car.ID,
		"marca":  input.Car.Marca,
		"modelo": input.Car.Modelo,
		"user_id": userID,
	})
	h.Publisher.Publish("car.created", meta)
	h.Logger.Info("car.created", "carro cadastrado", map[string]any{
		"car_id":  input.Car.ID,
		"user_id": userID,
	})

	c.JSON(http.StatusCreated, input.Car)
}

func (h *CarHandler) List(c *gin.Context) {
	cars, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cars)
}

func (h *CarHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	car, err := h.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "carro não encontrado"})
		return
	}
	c.JSON(http.StatusOK, car)
}

func (h *CarHandler) GetMyCars(c *gin.Context) {
	userID := c.GetInt("user_id")
	cars, err := h.Service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cars)
}

func (h *CarHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	var car models.Car
	if err := c.BindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	car.ID = id
	if err := h.Service.Update(&car); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, car)
}

func (h *CarHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
	if err := h.Service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	meta, _ := json.Marshal(map[string]any{"car_id": id})
	h.Publisher.Publish("car.deleted", meta)
	h.Logger.Info("car.deleted", "carro removido", map[string]any{"car_id": id})

	c.Status(http.StatusNoContent)
}
