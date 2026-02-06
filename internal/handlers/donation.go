package handlers

import (
	"net/http"
	"newapp/internal/models"
	"newapp/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DonationHandler struct {
	repo *repository.DonationRepository
}

func NewDonationHandler() *DonationHandler {
	return &DonationHandler{
		repo: repository.NewDonationRepository(),
	}
}

func (h *DonationHandler) Create(c *gin.Context) {
	var donation models.Donation
	if err := c.ShouldBindJSON(&donation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	donation.TempleID = 1 // Default temple

	if err := h.repo.Create(&donation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create donation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Donation recorded successfully",
		"data":    donation,
	})
}

func (h *DonationHandler) GetAll(c *gin.Context) {
	donations, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch donations"})
		return
	}

	total, _ := h.repo.GetTotal()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    donations,
		"total":   total,
	})
}

func (h *DonationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	donation, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Donation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    donation,
	})
}

func (h *DonationHandler) GetByFestival(c *gin.Context) {
	festivalID, err := strconv.ParseUint(c.Param("festivalId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid festival ID"})
		return
	}

	donations, err := h.repo.GetByFestival(uint(festivalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch donations"})
		return
	}

	total, _ := h.repo.GetTotalByFestival(uint(festivalID))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    donations,
		"total":   total,
	})
}

func (h *DonationHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var donation models.Donation
	if err := c.ShouldBindJSON(&donation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	donation.ID = uint(id)
	if err := h.repo.Update(&donation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update donation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Donation updated successfully",
		"data":    donation,
	})
}

func (h *DonationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete donation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Donation deleted successfully",
	})
}
