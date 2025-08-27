package services

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	if err := ensureIndexes(DB.Collection("posts")); err != nil {
		log.Fatal("❌ Error creando índices:", err)
	}
}


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