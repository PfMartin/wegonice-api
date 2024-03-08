package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID         string `json:"id" bson:"_id"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	Role       string `json:"role" bson:"role"`
	CreatedAt  int64  `json:"createdAt" bson:"createdAt"`
	ModifiedAt int64  `json:"modifiedAt" bson:"modifiedAt"`
}

type CollectionHandler struct {
	collection *mongo.Collection
}

func NewUsersHandler(dbClient *mongo.Client, dbName string) *CollectionHandler {
	collection := dbClient.Database(dbName).Collection("users")

	return &CollectionHandler{
		collection,
	}
}

func (handler *CollectionHandler) CreateUser(ctx context.Context, user User) (primitive.ObjectID, error) {
	insertData := bson.D{
		{Key: "email", Value: user.Email},
		{Key: "password", Value: user.Password},
		{Key: "role", Value: user.Role},
		{Key: "createdAt", Value: user.CreatedAt},
		{Key: "modifiedAt", Value: user.ModifiedAt},
	}

	cursor, err := handler.collection.InsertOne(ctx, insertData)
	if err != nil {
		return primitive.NilObjectID, err
	}

	userId := cursor.InsertedID.(primitive.ObjectID)

	return userId, nil
}
