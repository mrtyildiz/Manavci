package main

import (
	"github.com/gin-gonic/gin"
	"go-gin-api/config"
	"go-gin-api/routes"
)

func main() {
	r := gin.Default()
	config.ConnectDatabase()
	config.InitRedis()
	routes.SetupRoutes(r)
	r.Run(":5001")
}
