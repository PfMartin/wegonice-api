package db

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var sessionProjectStage = bson.M{"$project": bson.M{
	"_id":          1,
	"refreshToken": 1,
	"userAgent":    1,
	"isBlocked":    1,
	"clientIp":     1,
	"expiresAt":    1,
	"user": bson.M{
		"$arrayElemAt": bson.A{
			bson.M{"$map": bson.M{"input": "$user", "as": "user", "in": bson.M{
				"_id":   "$$user._id",
				"email": "$$user.email",
			},
			},
			}, 0,
		},
	},
}}

func (store *MongoDBStore) CreateSession(ctx context.Context, session Session) (primitive.ObjectID, error) {
	insertData := bson.M{
		"userId":       session.UserID,
		"refreshToken": session.RefreshToken,
		"userAgent":    session.UserAgent,
		"isBlocked":    session.IsBlocked,
		"clientIp":     session.ClientIP,
		"expiresAt":    session.ExpiresAt,
	}

	insertResult, err := store.sessionCollection.InsertOne(ctx, insertData)
	if err != nil {
		log.Err(err).Msg("failed to insert session")
		return primitive.NilObjectID, err
	}

	sessionID := insertResult.InsertedID.(primitive.ObjectID)

	return sessionID, nil
}

func (store *MongoDBStore) GetSessionByID(ctx context.Context, sessionID string) (Session, error) {
	var session Session

	primitiveSessionID, err := primitive.ObjectIDFromHex(sessionID)
	if err != nil {
		log.Err(err).Msgf("failed to parse sessionID %s to primitive ObjectID", sessionID)
		return session, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": primitiveSessionID}},
		userLookupStage,
		sessionProjectStage,
		{"$limit": 1},
	}

	cursor, err := store.sessionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Err(err).Msgf("failed to execute pipeline to find session with sessionID %s and its user", sessionID)
		return session, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		log.Error().Msgf("failed to find session with sessionID %s", sessionID)
		return session, fmt.Errorf("failed to find session with sessionID %s", sessionID)
	}

	if err := cursor.Decode(&session); err != nil {
		log.Err(err).Msg("failed to decode session")
		return session, nil
	}

	return session, nil
}
