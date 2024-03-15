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

	primitiveAuthorID, err := primitive.ObjectIDFromHex(recipe.AuthorID)
	if err != nil {
		log.Err(err).Msgf("failed to parse authorID %s to primitive ObjectID", recipe.AuthorID)
		return primitive.NilObjectID, err
	}

	primitiveUserID, err := primitive.ObjectIDFromHex(recipe.UserID)
	if err != nil {
		log.Err(err).Msgf("failed to parse userID %s to primitive ObjectID", recipe.UserID)
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
		"authorId":    primitiveAuthorID,
		"userId":      primitiveUserID,
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

// TODO: Aggregate Author and User into the recipes
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

// TODO: Aggregate author and user into the recipe
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

func (recipeColl *RecipeCollection) UpdateRecipeByID(ctx context.Context, recipeID string, recipeUpdate Recipe) (int64, error) {
	primitiveRecipeID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		log.Err(err).Msgf("failed to parse recipeID %s to primitive ObjectID", recipeID)
		return 0, err
	}

	filter := bson.M{
		"_id": primitiveRecipeID,
	}

	update := bson.M{
		"$set": bson.M{"modifiedAt": time.Now().Unix()},
	}
	// TODO: Create generic function for this
	if recipeUpdate.Name != "" {
		update["$set"].(bson.M)["name"] = recipeUpdate.Name
	}
	if recipeUpdate.ImageName != "" {
		update["$set"].(bson.M)["imageName"] = recipeUpdate.ImageName
	}
	if recipeUpdate.RecipeURL != "" {
		update["$set"].(bson.M)["recipeUrl"] = recipeUpdate.RecipeURL
	}
	if recipeUpdate.TimeM != 0 {
		update["$set"].(bson.M)["timeM"] = recipeUpdate.TimeM
	}
	if recipeUpdate.Category != "" {
		update["$set"].(bson.M)["category"] = recipeUpdate.Category
	}
	if len(recipeUpdate.Ingredients) != 0 {
		update["$set"].(bson.M)["ingredients"] = recipeUpdate.Ingredients
	}
	if len(recipeUpdate.PrepSteps) != 0 {
		update["$set"].(bson.M)["prepSteps"] = recipeUpdate.PrepSteps
	}
	if recipeUpdate.AuthorID != "" {
		update["$set"].(bson.M)["authorID"] = recipeUpdate.AuthorID
	}
	if recipeUpdate.UserID != "" {
		update["$set"].(bson.M)["userID"] = recipeUpdate.UserID
	}

	updateResult, err := recipeColl.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Err(err).Msgf("failed to update recipe with recipe recipeID %s", recipeID)
		return 0, err
	}

	if updateResult.MatchedCount < 1 {
		log.Info().Msgf("could not find recipe with recipeID %s", recipeID)
	}

	modifiedCount := updateResult.ModifiedCount
	if modifiedCount < 1 {
		log.Info().Msgf("did not update recipe with recipeID %s", recipeID)
	}

	return modifiedCount, err
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
