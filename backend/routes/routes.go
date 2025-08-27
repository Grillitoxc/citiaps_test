package routes

import (
	"github.com/gin-gonic/gin"
	"blog-api/controllers"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/posts", controllers.ListPosts)
		api.POST("/posts", controllers.CreatePost)
		api.GET("/posts/:id", controllers.GetPostByID)
		api.PUT("/posts/:id", controllers.UpdatePostByID)
		api.DELETE("/posts/:id", controllers.DeletePostByID)
		api.GET("/posts/metrics/by-tag", controllers.GetPostsMetricsByTag)
	}
}
