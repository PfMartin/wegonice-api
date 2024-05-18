package api

import (
	"net/http"

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
