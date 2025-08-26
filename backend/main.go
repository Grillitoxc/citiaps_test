package main

import (
    "log"

    "github.com/gin-gonic/gin"
	"blog-api/config"
    "blog-api/routes"
    "blog-api/services"
)

func main() {
	// Load config
  	cfg := config.Load()

	// Mongo DB connection
	services.ConnectMongo(cfg.MongoURI, cfg.MongoDB)

	// Router setup
	r := gin.Default()
	routes.SetupRoutes(r)

	// Start server
	if cfg.Port == "" {
		log.Fatal("❌ No se definió PORT en .env")
	}
	log.Println("🚀 API escuchando en :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}