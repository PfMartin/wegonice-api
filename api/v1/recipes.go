package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PfMartin/wegonice-api/db"
	"github.com/gin-gonic/gin"
)

// listRecipes
//
// @Summary			List all recipes
// @Description	All recipes are listed in a paginated manner
// @ID					recipes-list-recipes
// @Tags				recipes
// @Accept			json
// @Produce			json
// @Param				authorization	header			string							false	"Authorization header for bearer token"
// @Param				page_id				query 			int									true	"Offset for the pagination"
// @Param				page_size			query 			int									true	"Number of elements in one page"
// @Success			200						{array}			RecipeResponse						"List of recipes matching the given pagination parameters"
// @Failure			400						{object}		ErrorBadRequest						"Bad Request"
// @Failure			401						{object}		ErrorUnauthorized					"Unauthorized"
// @Failure 		500						{object}		ErrorInternalServerError	"Internal Server Error"
// @Router			/recipes			[get]
func (server *Server) listRecipes(ctx *gin.Context) {
	var pagination db.Pagination
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	// TODO: Add sorting

	recipes, err := server.store.GetAllRecipes(ctx, pagination)
	if err != nil {
		NewErrorInternalServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, recipes)
}

// createRecipe
//
// @Summary			Create new recipe
// @Description	Creates a new recipe
// @ID					recipes-create-recipe
// @Tags				recipes
// @Accept			json
// @Produce			json
// @Param				authorization		header			string							false	"Authorization header for bearer token"
// @Param				data						body 				RecipeToCreate			true	"Data for the recipe to create"
// @Success			201							string			string										"ID of the created recipe"
// @Failure			400							{object}		ErrorBadRequest						"Bad Request"
// @Failure			401							{object}		ErrorUnauthorized					"Unauthorized"
// @Failure 		500							{object}		ErrorInternalServerError	"Internal Server Error"
// @Router			/recipes				[post]
func (server *Server) createRecipe(ctx *gin.Context) {
	var recipeBody db.RecipeToCreate
	if err := ctx.ShouldBindJSON(&recipeBody); err != nil {
		fmt.Println(err)
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	recipeID, err := server.store.CreateRecipe(ctx, recipeBody)
	if err != nil {
		if strings.HasPrefix(err.Error(), "recipe with name") { // TODO: Find better way to check error types (enum?)
			NewErrorBadRequest(err).Send(ctx)
			return
		}

		NewErrorInternalServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, recipeID)
}

// getRecipeByID
//
// @Summary			Get one recipe by ID
// @Description	One recipe, which matches the ID, is returned
// @ID					recipes-get-recipe-by-id
// @Tags				recipes
// @Accept			json
// @Produce			json
// @Param				authorization		header			string							false	"Authorization header for bearer token"
// @Param				id							path 				int									true	"ID of the desired recipe"
// @Success			200							{object}		RecipeResponse						"Recipe that matches the ID"
// @Failure			400							{object}		ErrorBadRequest						"Bad Request"
// @Failure			401							{object}		ErrorUnauthorized					"Unauthorized"
// @Failure			404							{object}		ErrorNotFound							"Not Found"
// @Failure 		500							{object}		ErrorInternalServerError	"Internal Server Error"
// @Router			/recipes/{id}		[get]
func (server *Server) getRecipeByID(ctx *gin.Context) {
	var uriParam getByIDRequest
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	recipe, err := server.store.GetRecipeByID(ctx, uriParam.ID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "failed to find recipe") { // TODO: Find better method to distinguish between error types (enum?)
			NewErrorNotFound(err).Send(ctx)
			return
		}

		NewErrorInternalServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, recipe)
}
