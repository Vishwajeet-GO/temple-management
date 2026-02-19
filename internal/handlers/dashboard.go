package handlers

import (
	"net/http"

	"newapp/internal/database"
	"newapp/internal/models"

	"github.com/gin-gonic/gin"
)

func GetDashboardSummary(c *gin.Context) {
	var totalDon, totalExp float64
	var donC, expC, festC int64

	database.DB.Model(&models.Donation{}).Select("COALESCE(SUM(amount),0)").Scan(&totalDon)
	database.DB.Model(&models.Expense{}).Select("COALESCE(SUM(amount),0)").Scan(&totalExp)
	database.DB.Model(&models.Donation{}).Count(&donC)
	database.DB.Model(&models.Expense{}).Count(&expC)
	database.DB.Model(&models.Festival{}).Count(&festC)

	var recentDon []models.Donation
	var recentExp []models.Expense
	var upcoming []models.Festival
	database.DB.Preload("Festival").Order("created_at DESC").Limit(5).Find(&recentDon)
	database.DB.Preload("Festival").Order("created_at DESC").Limit(5).Find(&recentExp)
	database.DB.Where("status IN ?", []string{"upcoming", "ongoing"}).Order("start_date ASC").Limit(5).Find(&upcoming)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_income":       totalDon,
			"total_donations":    totalDon,
			"total_expenses":     totalExp,
			"balance":            totalDon - totalExp,
			"net_balance":        totalDon - totalExp,
			"donation_count":     donC,
			"expense_count":      expC,
			"festival_count":     festC,
			"recent_donations":   recentDon,
			"recent_expenses":    recentExp,
			"upcoming_festivals": upcoming,
		},
	})
}
