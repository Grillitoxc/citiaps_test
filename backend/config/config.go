// config/config.go
//
// Paquete config: carga y expone la configuración de la aplicación.
// 
// Convenciones:
//   - Se utiliza godotenv para cargar variables desde un archivo `.env` en desarrollo.
//   - En producción, se espera que las variables de entorno estén definidas en el sistema.
//   - El objeto Config centraliza los valores necesarios para arrancar el servidor
//     y la conexión a MongoDB.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config representa la configuración base de la aplicación.
//
// Campos:
//   - Port: puerto HTTP en el que se levanta la API (ej. ":8080").
//   - MongoURI: URI de conexión a MongoDB (ej. "mongodb://localhost:27017").
//   - MongoDB: nombre de la base de datos a utilizar.
type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
}

// Load inicializa la configuración cargando primero el archivo `.env` (si existe)
// y luego leyendo variables de entorno.
//
// Retorna:
//   - *Config con los valores de Port, MongoURI y MongoDB.
//
// Comportamiento:
//   - Si no encuentra un archivo `.env`, muestra un warning pero no falla.
//   - Las variables deben estar definidas en el entorno, ya sea por `.env` o por el sistema.
//   - En producción se recomienda definir las variables directamente en el sistema
//     y no depender de `.env`.
//
// Ejemplo de archivo .env:
//   PORT=:8080
//   MONGODB_URI=mongodb://localhost:27017
//   MONGODB_DB=blog
func Load() *Config {
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No se encontró archivo .env")
	}

	return &Config{
		Port:     getenv("PORT"),
		MongoURI: getenv("MONGODB_URI"),
		MongoDB:  getenv("MONGODB_DB"),
	}
}

// getenv retorna el valor de una variable de entorno.
//
// Parámetros:
//   - k: nombre de la variable de entorno.
//
// Retorna:
//   - string con el valor de la variable; "" si no existe.
func getenv(k string) string {
	return os.Getenv(k)
}
