package handlers

import (
	"net/http"
	"newapp/internal/models"
	"newapp/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	repo *repository.ExpenseRepository
}

func NewExpenseHandler() *ExpenseHandler {
	return &ExpenseHandler{
		repo: repository.NewExpenseRepository(),
	}
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.TempleID = 1 // Default temple

	if err := h.repo.Create(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Expense recorded successfully",
		"data":    expense,
	})
}

func (h *ExpenseHandler) GetAll(c *gin.Context) {
	expenses, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	total, _ := h.repo.GetTotal()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    expenses,
		"total":   total,
	})
}

func (h *ExpenseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	expense, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    expense,
	})
}

func (h *ExpenseHandler) GetByFestival(c *gin.Context) {
	festivalID, err := strconv.ParseUint(c.Param("festivalId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid festival ID"})
		return
	}

	expenses, err := h.repo.GetByFestival(uint(festivalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	total, _ := h.repo.GetTotalByFestival(uint(festivalID))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    expenses,
		"total":   total,
	})
}

func (h *ExpenseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense.ID = uint(id)
	if err := h.repo.Update(&expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Expense updated successfully",
		"data":    expense,
	})
}

func (h *ExpenseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Expense deleted successfully",
	})
}
