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
)

var authorLookupStage = bson.M{"$lookup": bson.M{
	"from":         "authors",
	"localField":   "authorId",
	"foreignField": "_id",
	"as":           "recipeAuthor",
}}

var authorProjectStage = bson.M{"$project": bson.M{
	"_id":          1,
	"name":         1,
	"firstName":    1,
	"lastName":     1,
	"websiteUrl":   1,
	"instagramUrl": 1,
	"youtubeUrl":   1,
	"imageName":    1,
	"createdAt":    1,
	"modifiedAt":   1,
	"userCreated": bson.M{
		"$arrayElemAt": bson.A{
			bson.M{"$map": bson.M{"input": "$user", "as": "userCreated", "in": bson.M{
				"_id":   "$$userCreated._id",
				"email": "$$userCreated.email",
			},
			},
			}, 0,
		},
	},
}}

func (store *MongoDBStore) CreateAuthor(ctx context.Context, author Author) (primitive.ObjectID, error) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := store.authorCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Err(err).Msgf("author with name %s already exists", author.Name)
		return primitive.NilObjectID, err
	}

	primitiveUserID, err := primitive.ObjectIDFromHex(author.UserID)
	if err != nil {
		log.Err(err).Msgf("failed to parse userID %s to primitive ObjectID", author.UserID)
		return primitive.NilObjectID, err
	}

	insertData := bson.M{
		"firstName":    author.FirstName,
		"lastName":     author.LastName,
		"name":         author.Name,
		"websiteUrl":   author.WebsiteURL,
		"instagramUrl": author.InstagramURL,
		"youtubeUrl":   author.YoutubeURL,
		"imageName":    author.ImageName,
		"userId":       primitiveUserID,
		"createdAt":    time.Now().Unix(),
		"modifiedAt":   time.Now().Unix(),
	}

	insertResult, err := store.authorCollection.InsertOne(ctx, insertData)
	if err != nil {
		log.Err(err).Msgf("failed to insert author with name %s", author.Name)
		return primitive.NilObjectID, err
	}

	authorID := insertResult.InsertedID.(primitive.ObjectID)

	return authorID, nil
}

func (store *MongoDBStore) GetAllAuthors(ctx context.Context, pagination Pagination) ([]Author, error) {
	var authors []Author

	pipeline := []bson.M{
		userLookupStage,
		authorProjectStage,
		getSortStage("name"),
		pagination.getSkipStage(),
		pagination.getLimitStage(),
	}

	cursor, err := store.authorCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Err(err).Msg("failed to aggregate author documents")
		return authors, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &authors); err != nil {
		log.Err(err).Msg("failed to parse author documents")
		return authors, err
	}

	return authors, nil
}

func (store *MongoDBStore) GetAuthorByID(ctx context.Context, authorID string) (Author, error) {
	var author Author

	primitiveAuthorID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		log.Err(err).Msgf("failed to parse authorID %s to primitive ObjectID", authorID)
		return author, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": primitiveAuthorID}},
		userLookupStage,
		authorProjectStage,
		{"$limit": 1},
	}

	cursor, err := store.authorCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Err(err).Msgf("failed to execute pipeline to find author with authorID %s and its user", authorID)
		return author, err
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		log.Error().Msgf("failed to find author with authorID %s", authorID)
		return author, fmt.Errorf("failed to find author with authorID %s", authorID)
	}

	if err := cursor.Decode(&author); err != nil {
		log.Err(err).Msg("failed to decode author")
		return author, nil
	}

	return author, nil
}

func (store *MongoDBStore) UpdateAuthorByID(ctx context.Context, authorID string, authorUpdate Author) (int64, error) {
	primitiveAuthorID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		log.Err(err).Msgf("failed to parse authorID %s to primitive ObjectID", authorID)
		return 0, err
	}

	filter := bson.M{
		"_id": primitiveAuthorID,
	}

	update := bson.M{
		"$set": bson.M{"modifiedAt": time.Now().Unix()},
	}
	if authorUpdate.Name != "" {
		update["$set"].(bson.M)["name"] = authorUpdate.Name
	}
	if authorUpdate.LastName != "" {
		update["$set"].(bson.M)["lastName"] = authorUpdate.LastName
	}
	if authorUpdate.FirstName != "" {
		update["$set"].(bson.M)["firstName"] = authorUpdate.FirstName
	}
	if authorUpdate.WebsiteURL != "" {
		update["$set"].(bson.M)["websiteUrl"] = authorUpdate.WebsiteURL
	}
	if authorUpdate.InstagramURL != "" {
		update["$set"].(bson.M)["instagramUrl"] = authorUpdate.InstagramURL
	}
	if authorUpdate.YoutubeURL != "" {
		update["$set"].(bson.M)["youtubeUrl"] = authorUpdate.YoutubeURL
	}
	if authorUpdate.ImageName != "" {
		update["$set"].(bson.M)["imageName"] = authorUpdate.ImageName
	}

	updateResult, err := store.authorCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Err(err).Msgf("failed to update author with author authorID %s", authorID)
		return 0, err
	}

	if updateResult.MatchedCount < 1 {
		log.Info().Msgf("could not find author with authorID %s", authorID)
	}

	modifiedCount := updateResult.ModifiedCount
	if modifiedCount < 1 {
		log.Info().Msgf("did not update author with authorID %s", authorID)
	}

	return modifiedCount, err
}

func (store *MongoDBStore) DeleteAuthorByID(ctx context.Context, authorID string) (int64, error) {
	primitiveAuthorID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		log.Err(err).Msgf("failed to parse authorID %s to primitive ObjectID", authorID)
		return 0, err
	}

	if err = checkReferencesOfDocument(ctx, store.recipeCollection, "authorId", primitiveAuthorID); err != nil {
		return 0, err
	}

	filter := bson.M{
		"_id": primitiveAuthorID,
	}

	deleteResult, err := store.authorCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Err(err).Msgf("failed to delete author with authorID %s", authorID)
		return 0, err
	}

	deleteCount := deleteResult.DeletedCount
	if deleteCount < 1 {
		log.Info().Msgf("author with authorID %s was not deleted", authorID)
	}

	return deleteCount, nil
}
