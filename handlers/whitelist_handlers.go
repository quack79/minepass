package handlers

import (
	"gabefraser/minepass/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorcon/rcon"
)

func WhitelistAdd(c *gin.Context) {
	rcon := c.MustGet("rcon").(*rcon.Conn)

	var req WhitelistUserRequest
	if err := c.BindJSON(&req); err != nil || !req.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No username was provided. Please try again.",
		})
		return
	}

	response, err := rcon.Execute("whitelist add " + req.Username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "Unable to update the whitelist"})
		return
	}

	utils.Logger(c.ClientIP() + " added " + req.Username + " to the whitelist")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": response,
	})
}

func WhitelistRemove(c *gin.Context) {
	rcon := c.MustGet("rcon").(*rcon.Conn)

	var req WhitelistUserRequest
	if err := c.BindJSON(&req); err != nil || !req.IsValid() {
		utils.Logger("No username was provided from " + c.ClientIP())

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No username was provided. Please try again.",
		})
		return
	}

	response, err := rcon.Execute("whitelist remove " + req.Username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "Unable to update the whitelist"})
		return
	}

	utils.Logger(c.ClientIP() + " removed " + req.Username + " from the whitelist")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": response,
	})
}

func WhitelistList(c *gin.Context) {
	rcon := c.MustGet("rcon").(*rcon.Conn)

	response, err := rcon.Execute("whitelist list")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"success": false, "message": "Unable to retrieve the whitelist"})
		return
	}

	players := parseWhitelist(response)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"players": players,
	})
}

func parseWhitelist(response string) []string {
	_, playerList, found := strings.Cut(response, ":")
	if !found || strings.TrimSpace(playerList) == "" {
		return []string{}
	}

	players := make([]string, 0)
	for _, player := range strings.Split(playerList, ",") {
		player = strings.TrimSpace(player)
		if minecraftUsername.MatchString(player) {
			players = append(players, player)
		}
	}

	return players
}
