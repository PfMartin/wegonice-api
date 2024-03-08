package db

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
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
	insertData := bson.M{
		"email":      user.Email,
		"password":   user.Password,
		"role":       user.Role,
		"createdAt":  user.CreatedAt,
		"modifiedAt": user.ModifiedAt,
	}

	cursor, err := handler.collection.InsertOne(ctx, insertData)
	if err != nil {
		return primitive.NilObjectID, err
	}

	userId := cursor.InsertedID.(primitive.ObjectID)

	return userId, nil
}

func (handler *CollectionHandler) GetAllUsers(ctx context.Context) ([]User, error) {
	var users []User

	cursor, err := handler.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Err(err).Msg("failed to find user documents")
		return users, err
	}

	if err = cursor.All(ctx, &users); err != nil {
		log.Err(err).Msg("failed to parse user documents")
		return users, err
	}

	return users, nil
}

func (handler *CollectionHandler) GetUserByID(ctx context.Context, userID string) (User, error) {
	var user User

	primitiveUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Err(err).Msgf("failed to parse userID %s to primitive ObjectID", userID)
		return user, err
	}

	filter := bson.M{
		"_id": primitiveUserID,
	}

	if err = handler.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		log.Err(err).Msgf("failed to find user with userID %s", userID)
		return user, nil
	}

	return user, nil
}

// TODO: Make this a patch method
func (handler *CollectionHandler) UpdateUserByID(ctx context.Context, userID string, userUpdate User) (int64, error) {
	primitiveUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Err(err).Msgf("failed to parse userID %s to primitive ObjectID", userID)
		return 0, err
	}

	filter := bson.M{
		"_id": primitiveUserID,
	}

	update := bson.M{
		"$set": bson.M{"modified_at": time.Now().Unix()},
	}
	if userUpdate.Email != "" {
		update["$set"].(bson.M)["email"] = userUpdate.Email
	}
	if userUpdate.Password != "" {
		update["$set"].(bson.M)["password"] = userUpdate.Password
	}

	updateResult, err := handler.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Err(err).Msgf("failed to update user with user userID %s", userID)
		return 0, err
	}

	if updateResult.MatchedCount < 1 {
		log.Info().Msgf("could not find user with userID %s", userID)
	}

	modifiedCount := updateResult.ModifiedCount
	if modifiedCount < 1 {
		log.Info().Msgf("did not update user with userID %s", userID)
	}

	return modifiedCount, err
}

func (handler *CollectionHandler) DeleteUserByID(ctx context.Context, userID string) (int64, error) {
	primitiveUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Err(err).Msgf("failed to parse userID %s to primitive ObjectID", userID)
		return 0, err
	}

	filter := bson.M{
		"_id": primitiveUserID,
	}

	deleteResult, err := handler.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Err(err).Msgf("failed to delete user with userID %s", userID)
		return 0, err
	}

	deleteCount := deleteResult.DeletedCount
	if deleteCount < 1 {
		log.Info().Msgf("user with userID %s was not deleted", userID)
	}

	return deleteCount, nil
}
