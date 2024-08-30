package main

import (
	"dirobash-api/configs"
	"dirobash-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	configs.ConnectDB()

	//routes
	routes.CitationRoute(router)

	router.Run(":" + configs.EnvPORT())
}
