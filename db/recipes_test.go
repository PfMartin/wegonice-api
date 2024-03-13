package db

import (
	"context"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/stretchr/testify/require"
)

var categories = []Category{Breakfast, Main, Desert, Smoothie, Baby, Drink}
var amountUnits = []AmountUnit{Milliliters, Liters, Milligrams, Grams, Tablespoon, Teaspoon, Piece}

func getRecipeCollection(t *testing.T) *RecipeCollection {
	conf := getDatabaseConfiguration(t)

	dbClient, _ := NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	require.NotNil(t, dbClient)

	coll := NewRecipeCollection(dbClient, conf.DBName)

	return coll
}

func createRandomRecipe(t *testing.T, recipeColl *RecipeCollection, userID string, authorID string) Recipe {
	var ingredients []Ingredient
	for i := 0; i < 10; i++ {
		amountIdx := util.RandomInt(0, int64(len(amountUnits)-1))

		ingredients = append(ingredients, Ingredient{
			Name:   util.RandomString(6),
			Amount: int(util.RandomInt(0, 100)),
			Unit:   amountUnits[amountIdx],
		})
	}

	var prepSteps []PrepStep
	for i := 0; i < 10; i++ {
		prepSteps = append(prepSteps, PrepStep{
			Rank:        i + 1,
			Description: util.RandomString(20),
		})
	}

	categoryIdx := util.RandomInt(0, int64(len(categories)-1))
	category := categories[categoryIdx]

	recipe := Recipe{
		Name:        util.RandomString(6),
		ImageName:   util.RandomString(10),
		RecipeURL:   util.RandomString(10),
		TimeM:       int(util.RandomInt(0, 180)),
		Category:    category,
		Ingredients: ingredients,
		PrepSteps:   prepSteps,
		AuthorID:    authorID,
		UserID:      userID,
	}

	insertedRecipeID, err := recipeColl.CreateRecipe(context.Background(), recipe)
	require.NoError(t, err)
	require.False(t, insertedRecipeID.IsZero())

	recipeID := insertedRecipeID.Hex()

	return Recipe{
		ID:          recipeID,
		Name:        recipe.Name,
		ImageName:   recipe.ImageName,
		RecipeURL:   recipe.RecipeURL,
		TimeM:       recipe.TimeM,
		Category:    recipe.Category,
		Ingredients: recipe.Ingredients,
		PrepSteps:   recipe.PrepSteps,
		AuthorID:    authorID,
		UserID:      userID,
		CreatedAt:   time.Now().Unix(),
		ModifiedAt:  time.Now().Unix(),
	}
}

func TestCreateRecipe(t *testing.T) {
	user := createRandomUser(t, getUserCollection(t))
	author := createRandomAuthor(t, getAuthorCollection(t), user.ID)

	recipeColl := getRecipeCollection(t)

	t.Run("Creates new recipe and throws an error when the same recipe should be created again", func(t *testing.T) {
		recipe := createRandomRecipe(t, recipeColl, user.ID, author.ID)

		_, err := recipeColl.CreateRecipe(context.Background(), recipe)
		require.Error(t, err)
	})
}

func TestGetAllRecipes(t *testing.T) {
	user := createRandomUser(t, getUserCollection(t))
	author := createRandomAuthor(t, getAuthorCollection(t), user.ID)

	recipeColl := getRecipeCollection(t)

	for i := 0; i < 10; i++ {
		_ = createRandomRecipe(t, recipeColl, user.ID, author.ID)
	}

	pagination := Pagination{
		PageID:   1,
		PageSize: 5,
	}

	t.Run("Gets all recipes with pagination", func(t *testing.T) {
		recipes, err := recipeColl.GetAllRecipes(context.Background(), pagination)
		require.NoError(t, err)
		require.NotEmpty(t, recipes)

		require.Equal(t, int(pagination.PageSize), len(recipes))
	})
}
