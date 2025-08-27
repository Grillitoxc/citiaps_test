// cmd/api/main.go
//
// Punto de entrada de la aplicación Blog API.
// Configura la aplicación cargando variables, inicializando la base de datos,
// configurando el router y levantando el servidor HTTP.
package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"blog-api/config"
	"blog-api/routes"
	"blog-api/services"
)

func main() {
	// 1. Cargar configuración (desde .env o entorno).
	cfg := config.Load()

	// 2. Conectar a MongoDB usando la configuración cargada.
	//    - Si la conexión o el ping fallan, el programa termina con log.Fatal.
	//    - Se crean índices necesarios en la colección "posts".
	services.ConnectMongo(cfg.MongoURI, cfg.MongoDB)

	// 3. Inicializar router con middlewares por defecto (logger + recovery).
	//    - gin.Default() incluye logging de requests y recuperación ante pánicos.
	r := gin.Default()

	// 4. Registrar las rutas de la API.
	//    - Ver routes.SetupRoutes: agrupa bajo /api y define endpoints de posts.
	routes.SetupRoutes(r)

	// 5. Iniciar servidor HTTP en el puerto configurado.
	if cfg.Port == "" {
		log.Fatal("❌ No se definió PORT en .env")
	}
	log.Println("🚀 API escuchando en :" + cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
