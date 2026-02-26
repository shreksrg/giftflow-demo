package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMe returns current user info
func GetMe(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}
