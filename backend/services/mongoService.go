// services/mongo.go
//
// Paquete services: inicialización de la conexión a MongoDB y configuración de índices.
//
// Convenciones:
//   - La variable global DB se inicializa al arrancar la aplicación mediante ConnectMongo.
//   - Se validan URI y nombre de la base de datos antes de intentar conectar.
//   - Se aplica un timeout de 10s en la conexión y ping.
//   - Al establecer la conexión, se crean índices obligatorios en la colección "posts".
package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB es la referencia global a la base de datos Mongo inicializada en ConnectMongo.
// Debe ser usada por los servicios para acceder a colecciones.
var DB *mongo.Database

// ConnectMongo establece la conexión a MongoDB y guarda la referencia global DB.
//
// Parámetros:
//   - uri: cadena de conexión a MongoDB (ej. "mongodb://localhost:27017").
//   - dbName: nombre de la base de datos a usar.
//
// Comportamiento:
//   - Valida que uri y dbName no sean cadenas vacías (termina la app con log.Fatal si lo son).
//   - Crea un cliente Mongo con timeout de 10s y realiza un Ping para verificar la conexión.
//   - Inicializa DB con la base especificada.
//   - Llama a ensureIndexes para asegurar la existencia de índices requeridos en la colección "posts".
//
// Errores:
//   - Si la conexión o el ping fallan, la aplicación termina con log.Fatal.
//   - Si la creación de índices falla, la aplicación también termina con log.Fatal.
func ConnectMongo(uri, dbName string) {
	if uri == "" || dbName == "" {
		log.Fatal("❌ Mongo config inválida (MONGODB_URI o MONGODB_DB vacíos)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Error conectando a Mongo:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Mongo no responde:", err)
	}

	DB = client.Database(dbName)
	log.Println("✅ Conectado a Mongo:", dbName)

	if err := ensureIndexes(DB.Collection("posts")); err != nil {
		log.Fatal("❌ Error creando índices:", err)
	}
}

// ensureIndexes crea los índices necesarios en la colección de posts.
//
// Índices definidos:
//   - text_title_content: índice de texto en {title, content} para búsquedas con $text.
//   - idx_published_publishedAt: índice compuesto en {published asc, publishedAt desc}
//     para optimizar listados paginados y filtros por estado de publicación.
//
// Retorna:
//   - error en caso de fallo en la creación de índices; nil si todo fue correcto.
func ensureIndexes(col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "title", Value: "text"},
				{Key: "content", Value: "text"},
			},
			Options: options.Index().SetName("text_title_content"),
		},
		{
			Keys: bson.D{
				{Key: "published", Value: 1},
				{Key: "publishedAt", Value: -1},
			},
			Options: options.Index().SetName("idx_published_publishedAt"),
		},
	})
	return err
}
