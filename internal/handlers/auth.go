package handlers

import (
	"net/http"
	"os"
	"time"

	"newapp/internal/database"
	"newapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password required"})
		return
	}

	var user models.User
	result := database.GetDB().Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is disabled"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(getJWTSecret()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", tokenString, 60*60*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, models.LoginResponse{
		Token:   tokenString,
		User:    user,
		Message: "Login successful",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var user models.User
	result := database.GetDB().First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	database.GetDB().First(&user, userID)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	database.GetDB().Model(&user).Update("password", string(hashedPassword))

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-key-change-in-production"
	}
	return secret
}
