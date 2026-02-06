package handlers

import (
	"net/http"
	"newapp/internal/database"
	"newapp/internal/models"
	"newapp/internal/repository"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	donationRepo *repository.DonationRepository
	expenseRepo  *repository.ExpenseRepository
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		donationRepo: repository.NewDonationRepository(),
		expenseRepo:  repository.NewExpenseRepository(),
	}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
	totalDonations, _ := h.donationRepo.GetTotal()
	totalExpenses, _ := h.expenseRepo.GetTotal()

	var totalFestivals, upcomingFestivals, totalEvents, donorsCount int64

	database.GetDB().Model(&models.Festival{}).Count(&totalFestivals)
	database.GetDB().Model(&models.Festival{}).Where("status = ?", "upcoming").Count(&upcomingFestivals)
	database.GetDB().Model(&models.Event{}).Count(&totalEvents)
	database.GetDB().Model(&models.Donation{}).Distinct("donor_phone").Count(&donorsCount)

	summary := models.DashboardSummary{
		TotalDonations:    totalDonations,
		TotalExpenses:     totalExpenses,
		Balance:           totalDonations - totalExpenses,
		TotalFestivals:    totalFestivals,
		UpcomingFestivals: upcomingFestivals,
		TotalEvents:       totalEvents,
		DonorsCount:       donorsCount,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

func (h *DashboardHandler) GetRecentDonations(c *gin.Context) {
	var donations []models.Donation
	database.GetDB().Order("created_at desc").Limit(10).Find(&donations)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    donations,
	})
}

func (h *DashboardHandler) GetRecentExpenses(c *gin.Context) {
	var expenses []models.Expense
	database.GetDB().Order("created_at desc").Limit(10).Find(&expenses)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    expenses,
	})
}
