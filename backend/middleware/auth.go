package middleware

import (
	"giftFlow/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Mock Users for simulation
var MockUsers = map[string]models.User{
	"dept_admin": {
		ID:         "user_001",
		Username:   "dept_admin",
		Role:       models.RoleDeptAdmin,
		Department: "Computer Science",
	},
	"dept_head": {
		ID:         "user_002",
		Username:   "dept_head",
		Role:       models.RoleDeptHead,
		Department: "Computer Science",
	},
	"cpro_admin": {
		ID:         "user_003",
		Username:   "cpro_admin",
		Role:       models.RoleCproAdmin,
		Department: "CPRO",
	},
}

// AuthMiddleware simulates OIDC authentication
// In production, this would validate a JWT token from Microsoft OIDC
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get username from header (simulating a login session or token claim)
		// For demo, we pass "X-Mock-User" header from frontend
		username := c.GetHeader("X-Mock-User")
		if username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Missing X-Mock-User header"})
			c.Abort()
			return
		}

		user, exists := MockUsers[username]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid mock user"})
			c.Abort()
			return
		}

		// Set user context for handlers
		c.Set("user", user)
		c.Next()
	}
}
