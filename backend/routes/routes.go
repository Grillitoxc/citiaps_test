// routes/routes.go
//
// Paquete routes: centraliza el enrutamiento HTTP de la API.
// Convenciones:
//   - Las rutas se agrupan bajo el prefijo "/api".
//   - No se implementa lógica aquí; únicamente se enrutan controladores.
package routes

import (
	"net/http"

	"blog-api/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes registra todas las rutas de la API en el router Gin.
//
// Endpoints principales:
//   - GET    /api/posts                  → listado con filtros y paginación
//   - POST   /api/posts                  → crear un post
//   - GET    /api/posts/:id              → obtener un post por ID
//   - PUT    /api/posts/:id              → actualizar un post por ID
//   - DELETE /api/posts/:id              → eliminar un post por ID
//   - GET    /api/posts/metrics/by-tag   → agregación: top-N tags por cantidad
//
// Adicionalmente, define un healthcheck en /healthz y manejadores
// para rutas no encontradas (404) y métodos no permitidos (405).
func SetupRoutes(r *gin.Engine) {
	// Healthcheck para test (Docker)
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.GET("/posts", controllers.ListPosts)
		api.POST("/posts", controllers.CreatePost)
		api.GET("/posts/:id", controllers.GetPostByID)
		api.PUT("/posts/:id", controllers.UpdatePostByID)
		api.DELETE("/posts/:id", controllers.DeletePostByID)
		api.GET("/posts/metrics/by-tag", controllers.GetPostsMetricsByTag)
	}

	// Handler global: 404 en JSON
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Ruta no encontrada",
			"path":    c.Request.URL.Path,
		})
	})

	// Handler global: 405 en JSON
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code":    http.StatusMethodNotAllowed,
			"message": "Método no permitido",
			"method":  c.Request.Method,
		})
	})
}
