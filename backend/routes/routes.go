package routes

import (
	"github.com/gin-gonic/gin"
	"blog-api/controllers"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/test", controllers.HelloWorld)
		
	}
}
