package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	WSHub          *Hub
	PacketsManager PacketsManagerI
)

func SetWebsocket() {
	WSHub = newWSHub()
	go WSHub.run()

	PacketsManager = newPacketsManager()
	go PacketsManager.HandleWSPackets()

	GamesManager := Games{
		Games: make(map[string]GameI),
	}

	go GamesManager.New()

	cf := CoinFlipGame{}
	cf.Register()

	API.GET("/gateway/websocket", func(c *gin.Context) {
		token := c.GetHeader("X-Auth-Token")
		if token == "" {
			JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		session := SessionTokens[token]
		if session.Token == "" || session.AccountID == "" {
			JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		user, found := GetUserByID(session.AccountID)
		if !found {
			JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		serveWs(WSHub, c.Writer, c.Request, user.ID)
	})
}
