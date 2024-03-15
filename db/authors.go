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

type AuthorCollection struct {
	collection *mongo.Collection
}

var authorProjection = bson.M{"$project": bson.M{
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
	"userId":       1,
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

	cursor, err := authorColl.collection.InsertOne(ctx, insertData)
	if err != nil {
		log.Err(err).Msgf("failed to insert author with name %s", author.Name)
		return primitive.NilObjectID, err
	}

	authorID := cursor.InsertedID.(primitive.ObjectID)

	return authorID, nil
}

// TODO: Aggregate user into the authors
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

// TODO: Aggregate user into the author
func (authorColl *AuthorCollection) GetAuthorByID(ctx context.Context, authorID string) (Author, error) {
	var author Author

	primitiveAuthorID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		log.Err(err).Msgf("failed to parse authorID %s to primitive ObjectID", authorID)
		return author, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": primitiveAuthorID}},
		userLookup,
		authorProjection,
		{"$limit": 1},
	}

	cursor, err := authorColl.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Err(err).Msg("failed to execute pipeline to find author and its user")
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
	recipeColl := NewRecipeCollection(authorColl.collection.Database().Client(), authorColl.collection.Database().Name())

	primitiveAuthorID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		log.Err(err).Msgf("failed to parse authorID %s to primitive ObjectID", authorID)
		return 0, err
	}

	count, err := recipeColl.collection.CountDocuments(ctx, bson.M{"authorId": primitiveAuthorID})
	if err != nil {
		return 0, err
	}
	if count > 0 {
		log.Error().Msg("can't delete author because it is referenced in at least one recipe.")
		return 0, fmt.Errorf("can't delete author because it is referenced in at least one recipe")
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
