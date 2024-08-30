package routes

import (
	"dirobash-api/controllers"

	"github.com/gin-gonic/gin"
)

func CitationRoute(router *gin.Engine) {
	router.POST("/citation", controllers.CreateCitation())
	router.GET("/citation/id/:citationId", controllers.GetACitationById())
	router.GET("/citation/number/:number", controllers.GetACitationByNumber())
	router.GET("/citations", controllers.GetAllCitations())
}
