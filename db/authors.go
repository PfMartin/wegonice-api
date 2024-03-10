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
		"firstName":   author.FirstName,
		"lastName":    author.LastName,
		"name":        author.Name,
		"website":     author.WebsiteURL,
		"instagram":   author.InstagramURL,
		"youTube":     author.YouTubeURL,
		"imageBase64": author.ImageBase64,
		"userId":      author.UserID,
		"createdAt":   time.Now().UnixMilli(),
		"modifiedAt":  time.Now().UnixMilli(),
	}

	cursor, err := authorColl.collection.InsertOne(ctx, insertData)
	if err != nil {
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
