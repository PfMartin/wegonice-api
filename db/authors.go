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

type AuthorCollection struct {
	collection *mongo.Collection
}

func NewAuthorCollection(dbClient *mongo.Client, dbName string) *AuthorCollection {
	collection := dbClient.Database(dbName).Collection("authors")

	return &AuthorCollection{
		collection,
	}
}

func (authorColl *AuthorCollection) CreateAuthor(ctx context.Context, author Author) (primitive.ObjectID, error) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := authorColl.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Err(err).Msgf("author with name %s already exists", author.Name)
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
		"userId":       author.UserID,
		"createdAt":    time.Now().Unix(),
		"modifiedAt":   time.Now().Unix(),
	}

	cursor, err := authorColl.collection.InsertOne(ctx, insertData)
	if err != nil {
		log.Err(err).Msgf("failed to insert author with name %s", author.Name)
		return primitive.NilObjectID, err
	}

	authorID := cursor.InsertedID.(primitive.ObjectID)

	return authorID, nil
}

func (authorColl *AuthorCollection) GetAllAuthors(ctx context.Context, pagination Pagination) ([]Author, error) {
	var authors []Author

	findOptions := pagination.getFindOptions()
	findOptions.SetSort(bson.M{"name": 1})

	cursor, err := authorColl.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Err(err).Msg("failed to find author documents")
		return authors, err
	}

	if err = cursor.All(ctx, &authors); err != nil {
		log.Err(err).Msg("failed to parse author documents")
		return authors, err
	}

	return authors, nil
}

func (authorColl *AuthorCollection) GetAuthorByID(ctx context.Context, authorID string) (Author, error) {
	var author Author

	primitiveAuthorID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		log.Err(err).Msgf("failed to parse authorID %s to primitive ObjectID", authorID)
		return author, err
	}

	filter := bson.M{
		"_id": primitiveAuthorID,
	}

	if err = authorColl.collection.FindOne(ctx, filter).Decode(&author); err != nil {
		log.Err(err).Msgf("failed to find author with authorID %s", authorID)
		return author, err
	}

	return author, nil
}

func (authorColl *AuthorCollection) UpdateAuthorByID(ctx context.Context, authorID string, authorUpdate Author) (int64, error) {
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

	updateResult, err := authorColl.collection.UpdateOne(ctx, filter, update)
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

func (authorColl *AuthorCollection) DeleteAuthorByID(ctx context.Context, authorID string) (int64, error) {
	primitiveAuthorID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		log.Err(err).Msgf("failed to parse authorID %s to primitive ObjectID", authorID)
		return 0, err
	}

	filter := bson.M{
		"_id": primitiveAuthorID,
	}

	deleteResult, err := authorColl.collection.DeleteOne(ctx, filter)
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