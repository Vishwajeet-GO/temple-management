package handlers

import (
    "net/http"
    "strconv"
    "newapp/internal/database"
    "newapp/internal/models"
    "github.com/gin-gonic/gin"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
    return &DashboardHandler{}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
    var totalDonations, totalExpenses float64
    var totalFestivals, upcomingFestivals int64

    database.GetDB().Model(&models.Donation{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalDonations)
    database.GetDB().Model(&models.Expense{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalExpenses)
    database.GetDB().Model(&models.Festival{}).Count(&totalFestivals)
    database.GetDB().Model(&models.Festival{}).Where("status = ?", "upcoming").Count(&upcomingFestivals)

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": models.DashboardSummary{
            TotalDonations:    totalDonations,
            TotalExpenses:     totalExpenses,
            Balance:           totalDonations - totalExpenses,
            TotalFestivals:    totalFestivals,
            UpcomingFestivals: upcomingFestivals,
        },
    })
}

// NEW: Get Stats for a Specific Project/Festival
func (h *DashboardHandler) GetProjectStats(c *gin.Context) {
    idStr := c.Param("id")
    id, _ := strconv.Atoi(idStr)

    var income, expense float64
    
    // Get total donations linked to this festival
    database.GetDB().Model(&models.Donation{}).Where("temple_id = ?", id).Select("COALESCE(SUM(amount), 0)").Scan(&income)
    
    // Note: In our current model, we are using the generic 'Purpose' string or we need to link IDs.
    // To make this robust, we assume the Frontend sends 'festival_id' in the request.
    // Let's query based on the Festival ID we added to models earlier.
    
    // We need to fix the query to check FestivalID column (assuming it exists in your DB from GORM auto-migrate)
    // If it doesn't exist, GORM will ignore it, but let's try to filter by the 'Purpose' text for now 
    // OR better, strict ID filtering.
    
    // Strict ID filtering (The correct way):
    database.GetDB().Model(&models.Donation{}).Where("festival_id = ?", id).Select("COALESCE(SUM(amount), 0)").Scan(&income)
    database.GetDB().Model(&models.Expense{}).Where("festival_id = ?", id).Select("COALESCE(SUM(amount), 0)").Scan(&expense)

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "income":  income,
        "expense": expense,
        "balance": income - expense,
    })
}

func (h *DashboardHandler) GetRecentDonations(c *gin.Context) {
    var donations []models.Donation
    database.GetDB().Order("created_at desc").Limit(5).Find(&donations)
    c.JSON(http.StatusOK, gin.H{"success": true, "data": donations})
}

func (h *DashboardHandler) GetRecentExpenses(c *gin.Context) {
    var expenses []models.Expense
    database.GetDB().Order("created_at desc").Limit(5).Find(&expenses)
    c.JSON(http.StatusOK, gin.H{"success": true, "data": expenses})
}
