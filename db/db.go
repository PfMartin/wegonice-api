package db

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDatabaseClient(authSource string, username string, password string, uri string) (*mongo.Client, context.CancelFunc) {
	credentials := options.Credential{
		AuthSource: authSource,
		Username:   username,
		Password:   password,
	}

	clientOptions := options.Client().ApplyURI(uri).SetAuth(credentials)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	dbClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		log.Fatal().Msgf("failed to connect to database: %s", err)
	}

	if err = dbClient.Ping(ctx, readpref.Primary()); err != nil {
		cancel()
		log.Fatal().Msgf("failed to ping database: %s", err)
	}

	log.Info().Msg("Connected to database")

	return dbClient, cancel
}

func getSortStage(key string) bson.M {
	return bson.M{"$sort": bson.M{key: 1}}
}

func checkReferencesOfDocument(ctx context.Context, coll *mongo.Collection, foreignKey string, id primitive.ObjectID) error {
	count, err := coll.CountDocuments(ctx, bson.M{foreignKey: id})
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("document with id %s is referenced in at least one other document", id)
	}

	return nil
}
