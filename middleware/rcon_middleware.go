package middleware

import (
	"gabefraser/minepass/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func RconMiddleware(c *gin.Context) {
	host := os.Getenv("MP_HOST")
	if host == "" {
		utils.Logger("Missing MP_HOST environment variable")
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "RCON host is not configured"})
		c.Abort()
		return
	}

	password := os.Getenv("MP_PASSWORD")
	if password == "" {
		utils.Logger("Missing MP_PASSWORD environment variable")
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "RCON password is not configured"})
		c.Abort()
		return
	}

	conn, err := rcon.Dial(host, password)
	if err != nil {
		utils.Logger("Unable to connect to RCON: " + err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "Unable to connect to the Minecraft RCON server"})
		c.Abort()
		return
	}
	defer conn.Close()

	c.Set("rcon", conn)
	c.Next()
}
