package middleware

import (
	"strings"
	"sync"

	"newapp/internal/models"

	"github.com/gin-gonic/gin"
)

var (
	Sessions = map[string]*models.User{}
	Mu       sync.RWMutex
)

func GetUser(token string) *models.User {
	Mu.RLock()
	defer Mu.RUnlock()
	return Sessions[token]
}

func SetUser(token string, user *models.User) {
	Mu.Lock()
	defer Mu.Unlock()
	Sessions[token] = user
}

func RemoveUser(token string) {
	Mu.Lock()
	defer Mu.Unlock()
	delete(Sessions, token)
}

func ExtractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return auth
}

func AuthRequired(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtractToken(c)
		user := GetUser(token)
		if user == nil {
			c.JSON(401, gin.H{"success": false, "error": "Login required"})
			c.Abort()
			return
		}
		if role == "admin" && user.Role != "admin" {
			c.JSON(403, gin.H{"success": false, "error": "Admin access required"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
