package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
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

func (recipeColl *RecipeCollection) CreateRecipe(ctx context.Context, recipe Recipe) {

}
