package main

import (
	"os"
	"time"
	"math/rand"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	Production = false
	API *gin.Engine
)

func main() {
	rand.Seed(time.Now().UnixNano())
	BoolGen = boolGen{}

	LoadEnviromentVariables()
	InitializeMongoConnection()
	SetHTTPServer()
	SetWebsocket()
	Listen()
}

func LoadEnviromentVariables() {
	envFilename := "dev.env"
	if Production {
		envFilename = "production.env"
	}

	godotenv.Load(envFilename)
}

func Listen() {
	API.Run(":" + os.Getenv("PORT"))
}