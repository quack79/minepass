package middleware

import (
	"gabefraser/minepass/utils"
	"os"

	"github.com/gin-gonic/gin"
)

func GuardMiddleware(c *gin.Context) {
	username := os.Getenv("MP_UI_USERNAME")
	if username == "" {
		username = "admin"
	}
	password := os.Getenv("MP_UI_PASSWORD")

	apiUsername := c.GetHeader("X-Api-Username")
	apiKey := c.GetHeader("X-Api-Key")
	if apiUsername == "" || apiKey == "" {
		utils.Logger("Missing credentials from the request from " + c.ClientIP())

		c.JSON(401, gin.H{
			"success": false,
			"message": "Missing username or password from the request",
		})
		c.Abort()
		return
	}

	if apiUsername != username || apiKey != password {
		utils.Logger("Incorrect credentials from " + c.ClientIP())

		c.JSON(401, gin.H{
			"success": false,
			"message": "Incorrect username or password",
		})
		c.Abort()
		return
	}

	c.Next()
}
