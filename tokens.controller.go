package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	SessionTokens = map[string]Session{}
)

func AssignToken(accountID string) (string, error) {
	token, err := GenerateToken(accountID)
	if err != nil {
		return "", err
	}

	SessionTokens[token] = Session{
		Token:     token,
		AccountID: accountID,
	}

	return token, nil
}

func RemoveToken(token string) {
	delete(SessionTokens, token)
}

func TokenMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Auth-Token")
		if token == "" {
			Abort(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		session := SessionTokens[token]
		if session.Token == "" || session.AccountID == "" {
			Abort(c, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		c.Set("accountID", session.AccountID)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-Auth-Token, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}