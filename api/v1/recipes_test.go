package api

import (
	"bytes"
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
	"github.com/gin-gonic/gin"
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
		{
			name:  "Success with pagination from 5 to 10",
			query: "?page_id=5&page_size=10",
			buildStubs: func(store *mock_db.MockDBStore) {
				pagination := db.Pagination{
					PageID:   5,
					PageSize: 10,
				}

				store.EXPECT().GetAllRecipes(gomock.Any(), pagination).Times(1).Return(recipes[4:], nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var gotRecipes []RecipeResponse
				err := json.NewDecoder(recorder.Body).Decode(&gotRecipes)
				require.NoError(t, err)

				require.Equal(t, 6, len(gotRecipes))

				for i, expectedRecipe := range recipes[4:] {
					requireRecipeComparison(t, expectedRecipe, gotRecipes[i])
				}
			},
		},
		{
			name:  "Fail with missing page_size",
			query: "?page_id=5",
			buildStubs: func(store *mock_db.MockDBStore) {
				pagination := db.Pagination{
					PageID: 5,
				}

				store.EXPECT().GetAllRecipes(gomock.Any(), pagination).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Fail with missing page_id",
			query: "?page_size=10",
			buildStubs: func(store *mock_db.MockDBStore) {
				pagination := db.Pagination{
					PageSize: 10,
				}

				store.EXPECT().GetAllRecipes(gomock.Any(), pagination).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

func TestUnitGetRecipeByID(t *testing.T) {
	user, _ := randomUser(t)
	var recipes []db.Recipe
	for i := 0; i < 2; i++ {
		recipe, _ := randomRecipe(t)
		recipes = append(recipes, recipe)
	}

	testCases := []struct {
		name          string
		id            string
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success getting the second recipe",
			id:   recipes[1].ID,
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().GetRecipeByID(gomock.Any(), recipes[1].ID).Times(1).Return(recipes[1], nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var gotRecipe RecipeResponse
				err := json.NewDecoder(recorder.Body).Decode(&gotRecipe)
				require.NoError(t, err)

				requireRecipeComparison(t, recipes[1], gotRecipe)
			},
		},
		{
			name: "Fail with non-existent ID",
			id:   "notexisting",
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().GetRecipeByID(gomock.Any(), "notexisting").Times(1).Return(db.Recipe{}, fmt.Errorf("failed to find recipe with recipeID: notexisting"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Fail with non-parsable ID",
			id:   "notexisting",
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().GetRecipeByID(gomock.Any(), "notexisting").Times(1).Return(db.Recipe{}, fmt.Errorf("failed to parse recipeID"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/api/v1/recipes/%s", tc.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUnitCreateRecipe(t *testing.T) {
	user, _ := randomUser(t)
	recipe, primitiveID := randomRecipe(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success creating a new recipe",
			body: gin.H{
				"name":        recipe.Name,
				"imageName":   recipe.ImageName,
				"recipeUrl":   recipe.RecipeURL,
				"timeM":       recipe.TimeM,
				"category":    recipe.Category,
				"ingredients": recipe.Ingredients,
				"prepSteps":   recipe.PrepSteps,
				"authorId":    recipe.AuthorID,
				"userId":      recipe.UserID,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateRecipe(gomock.Any(), db.RecipeToCreate{
					Name:        recipe.Name,
					ImageName:   recipe.ImageName,
					RecipeURL:   recipe.RecipeURL,
					TimeM:       recipe.TimeM,
					Category:    recipe.Category,
					Ingredients: recipe.Ingredients,
					PrepSteps:   recipe.PrepSteps,
					AuthorID:    recipe.AuthorID,
					UserID:      recipe.UserID,
				}).Times(1).Return(primitiveID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Fail due to missing name",
			body: gin.H{
				"imageName":   recipe.ImageName,
				"recipeUrl":   recipe.RecipeURL,
				"timeM":       recipe.TimeM,
				"category":    recipe.Category,
				"ingredients": recipe.Ingredients,
				"prepSteps":   recipe.PrepSteps,
				"authorId":    recipe.AuthorID,
				"userId":      recipe.UserID,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateRecipe(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fail due to missing userID",
			body: gin.H{
				"name":        recipe.Name,
				"imageName":   recipe.ImageName,
				"recipeUrl":   recipe.RecipeURL,
				"timeM":       recipe.TimeM,
				"category":    recipe.Category,
				"ingredients": recipe.Ingredients,
				"prepSteps":   recipe.PrepSteps,
				"authorId":    recipe.AuthorID,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateRecipe(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fail due to missing authorID",
			body: gin.H{
				"name":        recipe.Name,
				"imageName":   recipe.ImageName,
				"recipeUrl":   recipe.RecipeURL,
				"timeM":       recipe.TimeM,
				"category":    recipe.Category,
				"ingredients": recipe.Ingredients,
				"prepSteps":   recipe.PrepSteps,
				"userId":      recipe.UserID,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateRecipe(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fail due to already existing recipe",
			body: gin.H{
				"name":        recipe.Name,
				"imageName":   recipe.ImageName,
				"recipeUrl":   recipe.RecipeURL,
				"timeM":       recipe.TimeM,
				"category":    recipe.Category,
				"ingredients": recipe.Ingredients,
				"prepSteps":   recipe.PrepSteps,
				"authorId":    recipe.AuthorID,
				"userId":      recipe.UserID,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateRecipe(gomock.Any(), db.RecipeToCreate{
					Name:        recipe.Name,
					ImageName:   recipe.ImageName,
					RecipeURL:   recipe.RecipeURL,
					TimeM:       recipe.TimeM,
					Category:    recipe.Category,
					Ingredients: recipe.Ingredients,
					PrepSteps:   recipe.PrepSteps,
					AuthorID:    recipe.AuthorID,
					UserID:      recipe.UserID,
				}).Times(1).Return(primitive.NilObjectID, fmt.Errorf("recipe with name %s already exists", recipe.Name))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := "/api/v1/recipes/"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUnitPatchRecipeByID(t *testing.T) {
	user, _ := randomUser(t)
	recipe, _ := randomRecipe(t)
	nonMatchingID := primitive.NewObjectID().Hex()

	ingredients, prepSteps := randomIngredientsAndPrepSteps(t, 10, 10)

	fullRecipePatch := db.RecipeUpdate{
		Name:        util.RandomString(6),
		ImageName:   util.RandomString(6),
		RecipeURL:   util.RandomString(10),
		TimeM:       int(util.RandomInt(15, 180)),
		Category:    "breakfast",
		Ingredients: ingredients,
		PrepSteps:   prepSteps,
		AuthorID:    primitive.NewObjectID().Hex(),
	}

	testCases := []struct {
		name          string
		id            string
		body          gin.H
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success with full update of the recipe",
			id:   recipe.ID,
			body: gin.H{
				"name":        fullRecipePatch.Name,
				"imageName":   fullRecipePatch.ImageName,
				"recipeUrl":   fullRecipePatch.RecipeURL,
				"timeM":       fullRecipePatch.TimeM,
				"category":    fullRecipePatch.Category,
				"ingredients": fullRecipePatch.Ingredients,
				"prepSteps":   fullRecipePatch.PrepSteps,
				"authorId":    fullRecipePatch.AuthorID,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateRecipeByID(gomock.Any(), recipe.ID, fullRecipePatch).Times(1).Return(int64(1), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Success with partial update of the recipe",
			id:   recipe.ID,
			body: gin.H{
				"name": fullRecipePatch.Name,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateRecipeByID(gomock.Any(), recipe.ID, db.RecipeUpdate{
					Name: fullRecipePatch.Name,
				}).Times(1).Return(int64(1), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Fail due to missing id",
			id:   "",
			body: gin.H{
				"name": fullRecipePatch.Name,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateRecipeByID(gomock.Any(), recipe.ID, gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Fail due to missing body",
			id:   recipe.ID,
			body: gin.H{},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateRecipeByID(gomock.Any(), recipe.ID, gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fail due to the provided recipeID not being valid",
			id:   "not-valid-id",
			body: gin.H{
				"name": fullRecipePatch.Name,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateRecipeByID(gomock.Any(), "not-valid-id", db.RecipeUpdate{
					Name: fullRecipePatch.Name,
				}).Times(1).Return(int64(0), fmt.Errorf("failed to parse recipeID"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fail due to no matching recipe for recipe ID",
			id:   nonMatchingID,
			body: gin.H{
				"name": fullRecipePatch.Name,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateRecipeByID(gomock.Any(), nonMatchingID, db.RecipeUpdate{
					Name: fullRecipePatch.Name,
				}).Times(1).Return(int64(0), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := fmt.Sprintf("/api/v1/recipes/%s", tc.id)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireRecipeComparison(t *testing.T, expectedRecipe db.Recipe, gotRecipe RecipeResponse) {
	t.Helper()

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
