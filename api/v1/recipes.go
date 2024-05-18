package api

import (
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
