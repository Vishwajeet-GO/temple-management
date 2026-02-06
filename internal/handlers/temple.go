package handlers

import (
	"net/http"
	"newapp/internal/models"
	"newapp/internal/repository"

	"github.com/gin-gonic/gin"
)

type TempleHandler struct {
	repo *repository.TempleRepository
}

func NewTempleHandler() *TempleHandler {
	return &TempleHandler{
		repo: repository.NewTempleRepository(),
	}
}

func (h *TempleHandler) GetTemple(c *gin.Context) {
	temple, err := h.repo.GetTemple()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Temple not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    temple,
	})
}

func (h *TempleHandler) UpdateTemple(c *gin.Context) {
	var temple models.Temple
	if err := c.ShouldBindJSON(&temple); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.UpdateTemple(&temple); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update temple"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Temple updated successfully",
		"data":    temple,
	})
}
