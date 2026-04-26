package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	Service LogSvc
}

func (h *LogHandler) List(c *gin.Context) {
	level := c.Query("level")
	event := c.Query("event")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	logs, err := h.Service.GetAll(level, event, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
