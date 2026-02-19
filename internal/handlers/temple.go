package handlers

import (
	"net/http"

	"newapp/internal/database"
	"newapp/internal/models"

	"github.com/gin-gonic/gin"
)

func GetTemple(c *gin.Context) {
	var info models.TempleInfo
	database.DB.First(&info)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": info})
}

func UpdateTemple(c *gin.Context) {
	var info models.TempleInfo
	database.DB.First(&info)
	c.ShouldBindJSON(&info)
	database.DB.Save(&info)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": info})
}
