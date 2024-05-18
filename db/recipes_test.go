package db

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/stretchr/testify/require"
)

var categories = []Category{Breakfast, Main, Desert, Smoothie, Baby, Drink}
var amountUnits = []AmountUnit{Milliliters, Liters, Milligrams, Grams, Tablespoon, Teaspoon, Piece}

func getRandomIngredients(t *testing.T, ingredientsCount int, ingredients *[]Ingredient) {
	t.Helper()

	for i := 0; i < ingredientsCount; i++ {
		amountIdx := util.RandomInt(0, int64(len(amountUnits)-1))

		*ingredients = append(*ingredients, Ingredient{
			Name:   util.RandomString(6),
			Amount: int(util.RandomInt(0, 100)),
			Unit:   amountUnits[amountIdx],
		})
	}
}

func getRandomPrepSteps(t *testing.T, stepsCount int, prepSteps *[]PrepStep) {
	t.Helper()

	for i := 0; i < stepsCount; i++ {
		*prepSteps = append(*prepSteps, PrepStep{
			Rank:        i + 1,
			Description: util.RandomString(20),
		})
	}
}

func getRandomIngredientsAndPrepSteps(t *testing.T, ingredientCount, prepStepCount int) ([]Ingredient, []PrepStep) {
	t.Helper()

	var wg sync.WaitGroup

	var ingredients []Ingredient
	wg.Add(1)
	go func() {
		defer wg.Done()
		getRandomIngredients(t, ingredientCount, &ingredients)
	}()

	var prepSteps []PrepStep
	wg.Add(1)
	go func() {
		defer wg.Done()
		getRandomPrepSteps(t, prepStepCount, &prepSteps)
	}()

	wg.Wait()

	require.Equal(t, ingredientCount, len(ingredients))
	require.Equal(t, prepStepCount, len(prepSteps))

	return ingredients, prepSteps
}

func createRandomRecipe(t *testing.T, store *MongoDBStore, userID string, authorID string) Recipe {
	t.Helper()

	ingredients, prepSteps := getRandomIngredientsAndPrepSteps(t, 10, 10)

	categoryIdx := util.RandomInt(0, int64(len(categories)-1))
	category := categories[categoryIdx]

	recipe := RecipeToCreate{
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

	insertedRecipeID, err := store.CreateRecipe(context.Background(), recipe)
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

func TestUnitCreateRecipe(t *testing.T) {
	store := getMongoDBStore(t)

	user := createRandomUser(t, store)
	author := createRandomAuthor(t, store, user.ID)

	t.Run("Creates new recipe and throws an error when the same recipe should be created again", func(t *testing.T) {
		recipe := createRandomRecipe(t, store, user.ID, author.ID)

		recipeToCreate := RecipeToCreate{
			Name:        recipe.Name,
			ImageName:   recipe.ImageName,
			RecipeURL:   recipe.RecipeURL,
			TimeM:       recipe.TimeM,
			Category:    recipe.Category,
			Ingredients: recipe.Ingredients,
			PrepSteps:   recipe.PrepSteps,
			AuthorID:    recipe.AuthorID,
			UserID:      recipe.UserID,
		}

		_, err := store.CreateRecipe(context.Background(), recipeToCreate)
		require.Error(t, err)
	})
}

func TestUnitGetAllRecipes(t *testing.T) {
	store := getMongoDBStore(t)

	user := createRandomUser(t, store)
	author := createRandomAuthor(t, store, user.ID)

	for i := 0; i < 10; i++ {
		r := createRandomRecipe(t, store, user.ID, author.ID)
		fmt.Println(r.AuthorID)
	}

	pagination := Pagination{
		PageID:   1,
		PageSize: 5,
	}

	t.Run("Gets all recipes with pagination", func(t *testing.T) {
		ctx := context.Background()
		recipes, err := store.GetAllRecipes(ctx, pagination)

		for _, recipe := range recipes {
			fmt.Print(recipe.Name + " | ")
			fmt.Print(recipe.Author.Name + "\n")
		}

		require.NoError(t, err)
		require.NotEmpty(t, recipes)

		require.Equal(t, int(pagination.PageSize), len(recipes))

		for _, recipe := range recipes {
			require.NotEmpty(t, recipe.Author)
			require.NotEmpty(t, recipe.Author.Name)
			require.NotEmpty(t, recipe.Author.FirstName)
			require.NotEmpty(t, recipe.Author.LastName)
			require.NotEmpty(t, recipe.Author.WebsiteURL)
			require.NotEmpty(t, recipe.Author.YoutubeURL)
			require.NotEmpty(t, recipe.Author.ImageName)
			require.NotEmpty(t, recipe.Author.UserID)

			require.NotEmpty(t, recipe.UserCreated)
			require.NotEmpty(t, recipe.UserCreated.ID)
			require.NotEmpty(t, recipe.UserCreated.Email)

			require.Empty(t, recipe.UserID)
			require.Empty(t, recipe.AuthorID)
		}
	})
}

func TestUnitGetRecipeByID(t *testing.T) {
	store := getMongoDBStore(t)

	user := createRandomUser(t, store)
	author := createRandomAuthor(t, store, user.ID)

	createdRecipe := createRandomRecipe(t, store, user.ID, author.ID)

	testCases := []struct {
		name           string
		recipeID       string
		hasError       bool
		expectedRecipe Recipe
	}{
		{
			name:           "Success",
			recipeID:       createdRecipe.ID,
			hasError:       false,
			expectedRecipe: createdRecipe,
		},
		{
			name:     "Fail with invalid recipeID",
			recipeID: "test",
			hasError: true,
		},
		{
			name:     "Fail with recipeID not found",
			recipeID: "659c00751f717854f690270d",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotRecipe, err := store.GetRecipeByID(context.Background(), tc.recipeID)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			tc.expectedRecipe.UserCreated = User{
				ID:    user.ID,
				Email: user.Email,
			}

			tc.expectedRecipe.Author = Author{
				ID:           author.ID,
				Name:         author.Name,
				FirstName:    author.FirstName,
				LastName:     author.LastName,
				WebsiteURL:   author.WebsiteURL,
				InstagramURL: author.InstagramURL,
				YoutubeURL:   author.YoutubeURL,
				ImageName:    author.ImageName,
				UserID:       author.UserID,
			}

			tc.expectedRecipe.AuthorID = ""
			tc.expectedRecipe.UserID = ""

			require.NoError(t, err)
			require.Equal(t, tc.expectedRecipe, gotRecipe)
		})
	}
}

func TestUnitUpdateRecipeByID(t *testing.T) {
	store := getMongoDBStore(t)

	user := createRandomUser(t, store)
	author := createRandomAuthor(t, store, user.ID)

	createdRecipe := createRandomRecipe(t, store, user.ID, author.ID)

	ingredients, prepSteps := getRandomIngredientsAndPrepSteps(t, 5, 5)

	recipeUpdate := Recipe{
		Name:        util.RandomString(4),
		ImageName:   util.RandomString(10),
		RecipeURL:   util.RandomString(8),
		TimeM:       int(util.RandomInt(0, 180)),
		Category:    categories[util.RandomInt(0, int64(len(categories)-1))],
		Ingredients: ingredients,
		PrepSteps:   prepSteps,
	}

	testCases := []struct {
		name          string
		recipeID      string
		recipeUpdate  Recipe
		hasError      bool
		modifiedCount int64
	}{
		{
			name:          "Success",
			recipeID:      createdRecipe.ID,
			recipeUpdate:  recipeUpdate,
			hasError:      false,
			modifiedCount: 1,
		},
		{
			name:          "Fail with invalid recipeID",
			recipeID:      "test",
			hasError:      true,
			modifiedCount: 0,
		},
		{
			name:          "Fail with recipeID not found",
			recipeID:      "659c00751f717854f690270d",
			hasError:      false,
			modifiedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			modifiedCount, err := store.UpdateRecipeByID(context.Background(), tc.recipeID, tc.recipeUpdate)
			require.Equal(t, tc.modifiedCount, modifiedCount)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if modifiedCount < 1 {
				return
			}

			updatedRecipe, err := store.GetRecipeByID(context.Background(), tc.recipeID)
			require.NoError(t, err)

			expectedRecipe := Recipe{
				ID:          createdRecipe.ID,
				Name:        recipeUpdate.Name,
				ImageName:   recipeUpdate.ImageName,
				RecipeURL:   recipeUpdate.RecipeURL,
				TimeM:       recipeUpdate.TimeM,
				Category:    recipeUpdate.Category,
				Ingredients: recipeUpdate.Ingredients,
				PrepSteps:   recipeUpdate.PrepSteps,
				CreatedAt:   createdRecipe.CreatedAt,
				ModifiedAt:  time.Now().Unix(),
			}

			require.Equal(t, expectedRecipe.ID, updatedRecipe.ID)
			require.Equal(t, expectedRecipe.Name, updatedRecipe.Name)
			require.Equal(t, expectedRecipe.ImageName, updatedRecipe.ImageName)
			require.Equal(t, expectedRecipe.RecipeURL, updatedRecipe.RecipeURL)
			require.Equal(t, expectedRecipe.TimeM, updatedRecipe.TimeM)
			require.Equal(t, expectedRecipe.Category, updatedRecipe.Category)
			require.Equal(t, expectedRecipe.Ingredients, updatedRecipe.Ingredients)
			require.Equal(t, expectedRecipe.PrepSteps, updatedRecipe.PrepSteps)
			require.WithinDuration(t, time.Unix(expectedRecipe.CreatedAt, 0), time.Unix(updatedRecipe.CreatedAt, 0), 5*time.Second)
			require.WithinDuration(t, time.Unix(expectedRecipe.ModifiedAt, 0), time.Unix(updatedRecipe.ModifiedAt, 0), 5*time.Second)
		})
	}
}

func TestDeleteRecipeByID(t *testing.T) {
	store := getMongoDBStore(t)

	user := createRandomUser(t, store)
	author := createRandomAuthor(t, store, user.ID)

	createdRecipe := createRandomRecipe(t, store, user.ID, author.ID)

	testCases := []struct {
		name        string
		recipeID    string
		hasError    bool
		deleteCount int64
	}{
		{
			name:        "Success",
			recipeID:    createdRecipe.ID,
			hasError:    false,
			deleteCount: 1,
		},
		{
			name:        "Fail with invalid recipeID",
			recipeID:    "test",
			hasError:    true,
			deleteCount: 0,
		},
		{
			name:        "Fail with recipeID not found",
			recipeID:    "659c00751f717854f690270d",
			hasError:    false,
			deleteCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deleteCount, err := store.DeleteRecipeByID(context.Background(), tc.recipeID)
			require.Equal(t, tc.deleteCount, deleteCount)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tc.deleteCount, deleteCount)
		})
	}
}
