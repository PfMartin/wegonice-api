package db

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserCollection struct {
	collection *mongo.Collection
}

func NewUserCollection(dbClient *mongo.Client, dbName string) *UserCollection {
	collection := dbClient.Database(dbName).Collection("users")

	return &UserCollection{
		collection,
	}
}

func (handler *UserCollection) CreateUser(ctx context.Context, user User) (primitive.ObjectID, error) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := handler.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Err(err).Msgf("user with email %s already exists", user.Email)
		return primitive.NilObjectID, err
	}

	insertData := bson.M{
		"email":      user.Email,
		"password":   user.Password,
		"role":       user.Role,
		"createdAt":  time.Now().Unix(),
		"modifiedAt": time.Now().Unix(),
	}

	cursor, err := handler.collection.InsertOne(ctx, insertData)
	if err != nil {
		return primitive.NilObjectID, err
	}

	userID := cursor.InsertedID.(primitive.ObjectID)

	return userID, nil
}

func (handler *UserCollection) GetAllUsers(ctx context.Context, pagination Pagination) ([]User, error) {
	var users []User

	findOptions := pagination.getFindOptions()
	findOptions.SetSort(bson.M{"email": 1})

	cursor, err := handler.collection.Find(ctx, bson.M{}, findOptions)
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

func (handler *UserCollection) GetUserByID(ctx context.Context, userID string) (User, error) {
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
		return user, err
	}

	return user, nil
}

func (handler *UserCollection) UpdateUserByID(ctx context.Context, userID string, userUpdate User) (int64, error) {
	primitiveUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Err(err).Msgf("failed to parse userID %s to primitive ObjectID", userID)
		return 0, err
	}

	filter := bson.M{
		"_id": primitiveUserID,
	}

	update := bson.M{
		"$set": bson.M{"modifiedAt": time.Now().Unix()},
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

func (handler *UserCollection) DeleteUserByID(ctx context.Context, userID string) (int64, error) {
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
