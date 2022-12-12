package main

import (
	"github.com/gin-gonic/gin"
)

var (
	AuthenticationRouter *gin.RouterGroup
	v1                   *gin.RouterGroup
	GamesRouter          *gin.RouterGroup
)

func SetHTTPServer() {
	API = gin.Default()
	API.Use(CORSMiddleware())

	SetPathGroups()
}

func SetPathGroups() {
	v1 := API.Group("/v1")
	AuthenticationRouter = v1.Group("/auth")
	GamesRouter = v1.Group("/games")

	AuthenticationControllers()
}
