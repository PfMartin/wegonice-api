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

type RecipeCollection struct {
	collection *mongo.Collection
}

func NewRecipeCollection(dbClient *mongo.Client, dbName string) *RecipeCollection {
	collection := dbClient.Database(dbName).Collection("recipes")

	return &RecipeCollection{
		collection,
	}
}

func (recipeColl *RecipeCollection) CreateRecipe(ctx context.Context, recipe Recipe) (primitive.ObjectID, error) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := recipeColl.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Err(err).Msgf("recipe with name %s already exists", recipe.Name)
		return primitive.NilObjectID, err
	}

	insertData := bson.M{
		"name":        recipe.Name,
		"imageName":   recipe.ImageName,
		"recipeUrl":   recipe.RecipeURL,
		"timeM":       recipe.TimeM,
		"category":    recipe.Category,
		"ingredients": recipe.Ingredients,
		"prepSteps":   recipe.PrepSteps,
		"authorId":    recipe.AuthorID,
		"userId":      recipe.UserID,
		"createdAt":   time.Now().Unix(),
		"modifiedAt":  time.Now().Unix(),
	}

	cursor, err := recipeColl.collection.InsertOne(ctx, insertData)
	if err != nil {
		return primitive.NilObjectID, err
	}

	recipeID := cursor.InsertedID.(primitive.ObjectID)

	return recipeID, nil
}

func (recipeColl *RecipeCollection) GetAllRecipes(ctx context.Context, pagination Pagination) ([]Recipe, error) {
	var recipes []Recipe

	findOptions := pagination.getFindOptions()
	findOptions.SetSort(bson.M{"name": 1})

	cursor, err := recipeColl.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Err(err).Msg("failed to find recipe documents")
		return recipes, err
	}

	if err = cursor.All(ctx, &recipes); err != nil {
		log.Err(err).Msg("failed to parse recipe documents")
		return recipes, err
	}

	return recipes, nil
}

func (recipeColl *RecipeCollection) GetRecipeByID(ctx context.Context, recipeID string) (Recipe, error) {
	var recipe Recipe

	primitiveRecipeID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		log.Err(err).Msgf("failed to parse recipeID %s to primitive ObjectID", recipeID)
		return recipe, err
	}

	filter := bson.M{
		"_id": primitiveRecipeID,
	}

	if err = recipeColl.collection.FindOne(ctx, filter).Decode(&recipe); err != nil {
		log.Err(err).Msgf("failed to find recipe with recipeID %s", recipeID)
		return recipe, err
	}

	return recipe, nil
}

func (recipeColl *RecipeCollection) DeleteRecipeByID(ctx context.Context, recipeID string) (int64, error) {
	primitiveRecipeID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		log.Err(err).Msgf("failed to parse recipeID %s to primitive ObjectID", recipeID)
		return 0, err
	}

	filter := bson.M{
		"_id": primitiveRecipeID,
	}

	deleteResult, err := recipeColl.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Err(err).Msgf("failed to delete recipe with recipeID %s", recipeID)
		return 0, err
	}

	deleteCount := deleteResult.DeletedCount
	if deleteCount < 1 {
		log.Info().Msgf("recipe with recipeID %s was not deleted", recipeID)
	}

	return deleteCount, nil
}
