package handlers

import (
	"net/http"

	"newapp/internal/database"
	"newapp/internal/middleware"
	"newapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid data"})
		return
	}

	var user models.User
	if database.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user).Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid credentials"})
		return
	}

	token := uuid.New().String()
	middleware.SetUser(token, &user)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
		"user":    gin.H{"id": user.ID, "username": user.Username, "role": user.Role},
	})
}

func Logout(c *gin.Context) {
	token := middleware.ExtractToken(c)
	middleware.RemoveUser(token)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func AuthCheck(c *gin.Context) {
	token := middleware.ExtractToken(c)
	user := middleware.GetUser(token)
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user":          gin.H{"id": user.ID, "username": user.Username, "role": user.Role},
	})
}
