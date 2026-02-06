package handlers

import (
	"net/http"
	"newapp/internal/models"
	"newapp/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FestivalHandler struct {
	repo         *repository.FestivalRepository
	donationRepo *repository.DonationRepository
	expenseRepo  *repository.ExpenseRepository
}

func NewFestivalHandler() *FestivalHandler {
	return &FestivalHandler{
		repo:         repository.NewFestivalRepository(),
		donationRepo: repository.NewDonationRepository(),
		expenseRepo:  repository.NewExpenseRepository(),
	}
}

func (h *FestivalHandler) Create(c *gin.Context) {
	var festival models.Festival
	if err := c.ShouldBindJSON(&festival); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	festival.TempleID = 1 // Default temple

	if err := h.repo.Create(&festival); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create festival"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Festival created successfully",
		"data":    festival,
	})
}

func (h *FestivalHandler) GetAll(c *gin.Context) {
	festivals, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch festivals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    festivals,
	})
}

func (h *FestivalHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	festival, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Festival not found"})
		return
	}

	// Get totals
	totalDonations, _ := h.donationRepo.GetTotalByFestival(uint(id))
	totalExpenses, _ := h.expenseRepo.GetTotalByFestival(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    festival,
		"summary": gin.H{
			"total_donations": totalDonations,
			"total_expenses":  totalExpenses,
			"balance":         totalDonations - totalExpenses,
		},
	})
}

func (h *FestivalHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var festival models.Festival
	if err := c.ShouldBindJSON(&festival); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	festival.ID = uint(id)
	if err := h.repo.Update(&festival); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update festival"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Festival updated successfully",
		"data":    festival,
	})
}

func (h *FestivalHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete festival"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Festival deleted successfully",
	})
}

func (h *FestivalHandler) GetUpcoming(c *gin.Context) {
	festivals, err := h.repo.GetUpcoming()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch festivals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    festivals,
	})
}
