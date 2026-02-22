package handlers

import (
	"net/http"
	"time"

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

func GetAdminDashboard(c *gin.Context) {
	now := time.Now()
	currentMonth := now.Format("2006-01")
	startOfMonth := now.Format("2006-01") + "-01"
	endOfMonth := now.AddDate(0, 1, 0).Format("2006-01") + "-01"

	// Overall totals
	var totalDon, totalExp float64
	var donCount, expCount, festCount, galleryCount int64
	database.DB.Model(&models.Donation{}).Select("COALESCE(SUM(amount),0)").Scan(&totalDon)
	database.DB.Model(&models.Expense{}).Select("COALESCE(SUM(amount),0)").Scan(&totalExp)
	database.DB.Model(&models.Donation{}).Count(&donCount)
	database.DB.Model(&models.Expense{}).Count(&expCount)
	database.DB.Model(&models.Festival{}).Count(&festCount)
	database.DB.Model(&models.GalleryItem{}).Count(&galleryCount)

	// This month totals
	var monthDon, monthExp float64
	var monthDonCount, monthExpCount int64
	database.DB.Model(&models.Donation{}).Where("date >= ? AND date < ?", startOfMonth, endOfMonth).Select("COALESCE(SUM(amount),0)").Scan(&monthDon)
	database.DB.Model(&models.Expense{}).Where("date >= ? AND date < ?", startOfMonth, endOfMonth).Select("COALESCE(SUM(amount),0)").Scan(&monthExp)
	database.DB.Model(&models.Donation{}).Where("date >= ? AND date < ?", startOfMonth, endOfMonth).Count(&monthDonCount)
	database.DB.Model(&models.Expense{}).Where("date >= ? AND date < ?", startOfMonth, endOfMonth).Count(&monthExpCount)

	// Last 6 months data for chart
	type MonthData struct {
		Month     string  `json:"month"`
		Donations float64 `json:"donations"`
		Expenses  float64 `json:"expenses"`
	}
	var monthlyData []MonthData
	for i := 5; i >= 0; i-- {
		m := now.AddDate(0, -i, 0)
		mStart := m.Format("2006-01") + "-01"
		mEnd := m.AddDate(0, 1, 0).Format("2006-01") + "-01"
		var mDon, mExp float64
		database.DB.Model(&models.Donation{}).Where("date >= ? AND date < ?", mStart, mEnd).Select("COALESCE(SUM(amount),0)").Scan(&mDon)
		database.DB.Model(&models.Expense{}).Where("date >= ? AND date < ?", mStart, mEnd).Select("COALESCE(SUM(amount),0)").Scan(&mExp)
		monthlyData = append(monthlyData, MonthData{
			Month:     m.Format("Jan 2006"),
			Donations: mDon,
			Expenses:  mExp,
		})
	}

	// Top 5 donors
	type TopDonor struct {
		Donor string  `json:"donor"`
		Total float64 `json:"total"`
		Count int64   `json:"count"`
	}
	var topDonors []TopDonor
	database.DB.Model(&models.Donation{}).
		Select("donor, SUM(amount) as total, COUNT(*) as count").
		Group("donor").
		Order("total DESC").
		Limit(5).
		Scan(&topDonors)

	// Category-wise expenses
	type CategoryData struct {
		Category string  `json:"category"`
		Total    float64 `json:"total"`
		Count    int64   `json:"count"`
	}
	var categoryData []CategoryData
	database.DB.Model(&models.Expense{}).
		Select("category, SUM(amount) as total, COUNT(*) as count").
		Group("category").
		Order("total DESC").
		Scan(&categoryData)

	// Festival-wise breakdown
	type FestivalBreakdown struct {
		ID        uint    `json:"id"`
		Name      string  `json:"name"`
		Status    string  `json:"status"`
		Donations float64 `json:"donations"`
		Expenses  float64 `json:"expenses"`
		Balance   float64 `json:"balance"`
	}
	var festivals []models.Festival
	database.DB.Order("created_at DESC").Limit(10).Find(&festivals)
	var festBreakdown []FestivalBreakdown
	for _, f := range festivals {
		var fDon, fExp float64
		database.DB.Model(&models.Donation{}).Where("festival_id = ?", f.ID).Select("COALESCE(SUM(amount),0)").Scan(&fDon)
		database.DB.Model(&models.Expense{}).Where("festival_id = ?", f.ID).Select("COALESCE(SUM(amount),0)").Scan(&fExp)
		festBreakdown = append(festBreakdown, FestivalBreakdown{
			ID: f.ID, Name: f.Name, Status: f.Status,
			Donations: fDon, Expenses: fExp, Balance: fDon - fExp,
		})
	}

	// Recent activity (last 10 donations + expenses combined)
	var recentDon []models.Donation
	var recentExp []models.Expense
	database.DB.Preload("Festival").Order("created_at DESC").Limit(10).Find(&recentDon)
	database.DB.Preload("Festival").Order("created_at DESC").Limit(10).Find(&recentExp)

	// Pending counts
	var pendingDon, pendingExp int64
	database.DB.Model(&models.Donation{}).Where("status = ?", "pending").Count(&pendingDon)
	database.DB.Model(&models.Expense{}).Where("status = ?", "pending").Count(&pendingExp)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"current_month":      currentMonth,
			"total_donations":    totalDon,
			"total_expenses":     totalExp,
			"balance":            totalDon - totalExp,
			"donation_count":     donCount,
			"expense_count":      expCount,
			"festival_count":     festCount,
			"gallery_count":      galleryCount,
			"month_donations":    monthDon,
			"month_expenses":     monthExp,
			"month_don_count":    monthDonCount,
			"month_exp_count":    monthExpCount,
			"month_balance":      monthDon - monthExp,
			"monthly_chart":      monthlyData,
			"top_donors":         topDonors,
			"category_expenses":  categoryData,
			"festival_breakdown": festBreakdown,
			"recent_donations":   recentDon,
			"recent_expenses":    recentExp,
			"pending_donations":  pendingDon,
			"pending_expenses":   pendingExp,
		},
	})
}

func UpdatePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid data"})
		return
	}

	user, _ := c.Get("user")
	u := user.(*models.User)

	// Verify old password
	var dbUser models.User
	if database.DB.Where("id = ? AND password = ?", u.ID, req.OldPassword).First(&dbUser).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Current password is incorrect"})
		return
	}

	if len(req.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Password must be at least 6 characters"})
		return
	}

	database.DB.Model(&dbUser).Update("password", req.NewPassword)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Password updated successfully"})
}
