package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDatabase(authSource string, username string, password string, uri string) *mongo.Client {
	credentials := options.Credential{
		AuthSource: authSource,
		Username:   username,
		Password:   password,
	}

	fmt.Println(authSource, username, password, uri)

	clientOptions := options.Client().ApplyURI(uri).SetAuth(credentials)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	dbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err = dbClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")
	return dbClient
}
