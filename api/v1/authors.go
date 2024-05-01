package api

import (
	"net/http"

	"github.com/PfMartin/wegonice-api/db"
	"github.com/gin-gonic/gin"
)

// listAuthors
//
// @Summary			List all authors
// @Description	All authors are listed in a paginated manner
// @ID					authors-list-authors
// @Tags				authors
// @Accept			json
// @Produce			json
// @Param				page_id			int					true											"Offset for the pagination"
// @Param				page_size		int					true											"Number of elements in one page"
// @Success			200					{array}			author										"List of authors matching the given pagination parameters"
// @Failure			400					{object}		ErrorBadRequest						"Bad Request"
// @Failure			401					{object}		ErrorUnauthorized					"Unauthorized"
// @Failure 		500					{object}		ErrorInternalServerError	"Internal Server Error"
// @Router			/authors		[get]
func (server *Server) listAuthors(ctx *gin.Context) {
	var pagination db.Pagination
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	authors, err := server.store.GetAllAuthors(ctx, pagination)
	if err != nil {
		NewErrorInternalServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, authors)
}
