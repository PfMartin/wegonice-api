package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/db"
	mock_db "github.com/PfMartin/wegonice-api/db/mock"
	"github.com/PfMartin/wegonice-api/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var amountUnits = []db.AmountUnit{db.Milliliters, db.Liters, db.Milligrams, db.Grams, db.Tablespoon, db.Teaspoon, db.Piece}

func randomIngredients(t *testing.T, ingredientsCount int, ingredients *[]db.Ingredient) {
	t.Helper()

	for i := 0; i < ingredientsCount; i++ {
		amountIdx := util.RandomInt(0, int64(len(amountUnits)-1))

		*ingredients = append(*ingredients, db.Ingredient{
			Name:   util.RandomString(6),
			Amount: int(util.RandomInt(0, 100)),
			Unit:   amountUnits[amountIdx],
		})
	}
}

func randomPrepSteps(t *testing.T, stepsCount int, prepSteps *[]db.PrepStep) {
	t.Helper()

	for i := 0; i < stepsCount; i++ {
		*prepSteps = append(*prepSteps, db.PrepStep{
			Rank:        i + 1,
			Description: util.RandomString(20),
		})
	}
}

func randomIngredientsAndPrepSteps(t *testing.T, ingredientCount, prepStepCount int) ([]db.Ingredient, []db.PrepStep) {
	t.Helper()

	var wg sync.WaitGroup

	var ingredients []db.Ingredient
	wg.Add(1)
	go func() {
		defer wg.Done()
		randomIngredients(t, ingredientCount, &ingredients)
	}()

	var prepSteps []db.PrepStep
	wg.Add(1)
	go func() {
		defer wg.Done()
		randomPrepSteps(t, prepStepCount, &prepSteps)
	}()

	wg.Wait()

	require.Equal(t, ingredientCount, len(ingredients))
	require.Equal(t, prepStepCount, len(prepSteps))

	return ingredients, prepSteps
}

func randomRecipe(t *testing.T) (db.Recipe, primitive.ObjectID) {
	t.Helper()

	userID := primitive.NewObjectID().Hex()
	recipeID := primitive.NewObjectID()

	author, authorID := randomAuthor(t)

	ingredients, prepSteps := randomIngredientsAndPrepSteps(t, 5, 10)

	return db.Recipe{
		ID:        recipeID.Hex(),
		Name:      util.RandomString(6),
		ImageName: util.RandomString(6),
		RecipeURL: util.RandomString(10),
		TimeM:     int(util.RandomInt(5, 120)),
		Category:  db.Breakfast,
		AuthorID:  authorID.Hex(),
		UserID:    userID,
		UserCreated: db.User{
			ID:    userID,
			Email: util.RandomEmail(),
		},
		Author:      author,
		Ingredients: ingredients,
		PrepSteps:   prepSteps,
	}, recipeID
}

func TestUnitListRecipes(t *testing.T) {
	user, _ := randomUser(t)
	var recipes []db.Recipe
	for i := 0; i < 10; i++ {
		recipe, _ := randomRecipe(t)
		recipes = append(recipes, recipe)
	}

	testCases := []struct {
		name          string
		query         string
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "Success with pagination from 1 to 10",
			query: "?page_id=1&page_size=10",
			buildStubs: func(store *mock_db.MockDBStore) {
				pagination := db.Pagination{
					PageID:   1,
					PageSize: 10,
				}

				store.EXPECT().GetAllRecipes(gomock.Any(), pagination).Times(1).Return(recipes, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var gotRecipes []RecipeResponse
				err := json.NewDecoder(recorder.Body).Decode(&gotRecipes)
				require.NoError(t, err)

				require.Equal(t, 10, len(gotRecipes))

				for i, expectedRecipe := range recipes {
					requireRecipeComparison(t, expectedRecipe, gotRecipes[i])
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockDBStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/recipes%s", tc.query)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireRecipeComparison(t *testing.T, expectedRecipe db.Recipe, gotRecipe RecipeResponse) {
	require.Equal(t, expectedRecipe.ID, gotRecipe.ID)
	require.Equal(t, expectedRecipe.Name, gotRecipe.Name)
	require.Equal(t, expectedRecipe.ImageName, gotRecipe.ImageName)
	require.Equal(t, expectedRecipe.Category, gotRecipe.Category)
	require.Equal(t, expectedRecipe.TimeM, gotRecipe.TimeM)
	require.Equal(t, expectedRecipe.RecipeURL, gotRecipe.RecipeURL)
	require.Equal(t, expectedRecipe.AuthorID, gotRecipe.AuthorID)
	require.Equal(t, expectedRecipe.UserID, gotRecipe.UserID)
	require.Equal(t, expectedRecipe.Ingredients, gotRecipe.Ingredients)
	require.Equal(t, expectedRecipe.PrepSteps, gotRecipe.PrepSteps)

	requireAuthorComparison(t, expectedRecipe.Author, gotRecipe.Author)
	requireUserComparison(t, expectedRecipe.UserCreated, gotRecipe.UserCreated)
}
