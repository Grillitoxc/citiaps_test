package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

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
}
