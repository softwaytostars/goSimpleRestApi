package main

import (
	"github.com/gin-gonic/gin"
	"goapi/resources"
	"goapi/services"
	"os"
)

func configureRouter() *gin.Engine {
	router := gin.Default()
	resources.RegisterHandlers(router, services.NewDocumentServiceImpl())
	return router
}

func main() {
	router := configureRouter()
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8040"
	}
	err := router.Run(":" + httpPort)
	if err != nil {
		return
	}
}
