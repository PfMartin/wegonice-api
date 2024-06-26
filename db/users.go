package db

import (
	"context"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userLookupStage = bson.M{"$lookup": bson.M{
	"from":         "users",
	"localField":   "userId",
	"foreignField": "_id",
	"as":           "user",
}}

func (store *MongoDBStore) CreateUser(ctx context.Context, user User) (primitive.ObjectID, error) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := store.userCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Err(err).Msgf("user with email %s already exists", user.Email)
		return primitive.NilObjectID, err
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		log.Err(err).Msgf("failed to hash password")
		return primitive.NilObjectID, err
	}

	if user.Role == "" {
		user.Role = "user"
	}

	insertData := bson.M{
		"email":        user.Email,
		"passwordHash": hashedPassword,
		"role":         user.Role,
		"isActive":     user.IsActive,
		"createdAt":    time.Now().Unix(),
		"modifiedAt":   time.Now().Unix(),
	}

	insertResult, err := store.userCollection.InsertOne(ctx, insertData)
	if err != nil {
		log.Err(err).Msgf("failed to insert user with email %s", user.Email)
		return primitive.NilObjectID, err
	}

	userID := insertResult.InsertedID.(primitive.ObjectID)

	return userID, nil
}

func (store *MongoDBStore) GetAllUsers(ctx context.Context, pagination Pagination) ([]User, error) {
	var users []User

	findOptions := pagination.getFindOptions()
	findOptions.SetSort(bson.M{"email": 1})

	cursor, err := store.userCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Err(err).Msg("failed to find user documents")
		return users, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		log.Err(err).Msg("failed to parse user documents")
		return users, err
	}

	return users, nil
}

func (store *MongoDBStore) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var user User

	filter := bson.M{
		"email": email,
	}

	if err := store.userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		log.Err(err).Msgf("failed to find user with email %s", email)
		return user, err
	}

	return user, nil
}

func (store *MongoDBStore) GetUserByID(ctx context.Context, userID string) (User, error) {
	var user User

	primitiveUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Err(err).Msgf("failed to parse userID %s to primitive ObjectID", userID)
		return user, err
	}

	filter := bson.M{
		"_id": primitiveUserID,
	}

	if err = store.userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		log.Err(err).Msgf("failed to find user with userID %s", userID)
		return user, err
	}

	return user, nil
}

func (store *MongoDBStore) UpdateUserByID(ctx context.Context, userID string, userUpdate User) (int64, error) {
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
		hashedPassword, err := util.HashPassword(userUpdate.Password)
		if err != nil {
			log.Err(err).Msg("failed to hash new password")
			return 0, err
		}
		update["$set"].(bson.M)["passwordHash"] = hashedPassword
	}

	updateResult, err := store.userCollection.UpdateOne(ctx, filter, update)
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

func (store *MongoDBStore) DeleteUserByID(ctx context.Context, userID string) (int64, error) {
	primitiveUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Err(err).Msgf("failed to parse userID %s to primitive ObjectID", userID)
		return 0, err
	}

	if err = checkReferencesOfDocument(ctx, store.recipeCollection, "userId", primitiveUserID); err != nil {
		return 0, err
	}

	if err = checkReferencesOfDocument(ctx, store.authorCollection, "userId", primitiveUserID); err != nil {
		return 0, err
	}

	filter := bson.M{
		"_id": primitiveUserID,
	}

	deleteResult, err := store.userCollection.DeleteOne(ctx, filter)
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
