package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo *mongo.Database

func ConnectMongo(uri, dbname string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("mongo connect error:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("mongo ping error:", err)
	}

	Mongo = client.Database(dbname)
}
