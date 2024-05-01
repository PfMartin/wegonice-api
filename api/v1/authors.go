package api

import (
	"net/http"
	"strings"

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
// @Param				page_id			path 				int									true	"Offset for the pagination"
// @Param				page_size		path 				int									true	"Number of elements in one page"
// @Success			200					{array}			AuthorResponse						"List of authors matching the given pagination parameters"
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

	// TODO: Add sorting

	authors, err := server.store.GetAllAuthors(ctx, pagination)
	if err != nil {
		NewErrorInternalServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, authors)
}

// getAuthorByID
//
// @Summary			Get one author by ID
// @Description	One author, which matches the ID is returned
// @ID					authors-get-author-by-id
// @Tags				authors
// @Accept			json
// @Produce			json
// @Param				id							path 				int									true	"ID of the desired author"
// @Success			200							{object}		AuthorResponse						"Author that matches the ID"
// @Failure			400							{object}		ErrorBadRequest						"Bad Request"
// @Failure			401							{object}		ErrorUnauthorized					"Unauthorized"
// @Failure			404							{object}		ErrorNotFound							"Not Found"
// @Failure 		500							{object}		ErrorInternalServerError	"Internal Server Error"
// @Router			/authors/{id}		[get]
func (server *Server) getAuthorByID(ctx *gin.Context) {
	var uriParam getByIDRequest
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	author, err := server.store.GetAuthorByID(ctx, uriParam.ID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "failed to find author") { // TODO: Find better method to distinguish between error types (enum?)
			NewErrorNotFound(err).Send(ctx)
			return
		}

		NewErrorInternalServerError(err).Send(ctx)
	}

	ctx.JSON(http.StatusOK, author)
}
