package db

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDatabase(authSource string, username string, password string, uri string) (*mongo.Client, context.CancelFunc) {
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
