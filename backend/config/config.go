package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No se encontró archivo .env")
	}

	cfg := &Config{
		Port:     getenv("PORT"),
		MongoURI: getenv("MONGODB_URI"),
		MongoDB:  getenv("MONGODB_DB"),
	}

	return cfg
}

func getenv(k string) string {
	return os.Getenv(k)
}
