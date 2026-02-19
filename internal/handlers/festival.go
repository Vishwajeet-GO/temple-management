package handlers

import (
	"net/http"
	"newapp/internal/database"
	"newapp/internal/models"

	"github.com/gin-gonic/gin"
)

func GetFestivals(c *gin.Context) {
	var festivals []models.Festival
	status := c.Query("status")
	query := database.DB.Order("created_at DESC")
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	query.Find(&festivals)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": festivals})
}

func GetFestivalReport(c *gin.Context) {
	id := c.Param("id")
	var festival models.Festival
	if database.DB.First(&festival, id).Error != nil {
		c.JSON(404, gin.H{"success": false, "error": "Festival not found"})
		return
	}

	var totalDon, totalExp float64
	var donCount, expCount int64
	database.DB.Model(&models.Donation{}).Where("festival_id=?", id).Select("COALESCE(SUM(amount),0)").Scan(&totalDon)
	database.DB.Model(&models.Expense{}).Where("festival_id=?", id).Select("COALESCE(SUM(amount),0)").Scan(&totalExp)
	database.DB.Model(&models.Donation{}).Where("festival_id=?", id).Count(&donCount)
	database.DB.Model(&models.Expense{}).Where("festival_id=?", id).Count(&expCount)

	var donations []models.Donation
	var expenses []models.Expense
	database.DB.Where("festival_id=?", id).Order("created_at DESC").Find(&donations)
	database.DB.Where("festival_id=?", id).Order("created_at DESC").Find(&expenses)

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"festival":        festival,
			"total_donations": totalDon,
			"total_expenses":  totalExp,
			"balance":         totalDon - totalExp,
			"donation_count":  donCount,
			"expense_count":   expCount,
			"donations":       donations,
			"expenses":        expenses,
		},
	})
}

func CreateFestival(c *gin.Context) {
	var f models.Festival
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(400, gin.H{"success": false, "error": "Invalid data: " + err.Error()})
		return
	}
	if f.Name == "" {
		c.JSON(400, gin.H{"success": false, "error": "Festival name required"})
		return
	}
	if f.Status == "" {
		f.Status = "upcoming"
	}
	database.DB.Create(&f)
	c.JSON(201, gin.H{"success": true, "data": f, "message": "Festival added"})
}

func UpdateFestival(c *gin.Context) {
	id := c.Param("id")
	var f models.Festival
	if database.DB.First(&f, id).Error != nil {
		c.JSON(404, gin.H{"success": false, "error": "Not found"})
		return
	}
	var u models.Festival
	c.ShouldBindJSON(&u)
	database.DB.Model(&f).Updates(u)
	c.JSON(200, gin.H{"success": true, "data": f, "message": "Updated"})
}

func DeleteFestival(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Festival{}, id)
	c.JSON(200, gin.H{"success": true, "message": "Deleted"})
}
