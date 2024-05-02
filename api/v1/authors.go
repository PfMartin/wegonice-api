package api

import (
	"fmt"
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
// @Param				page_id			query 			int									true	"Offset for the pagination"
// @Param				page_size		query 			int									true	"Number of elements in one page"
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

// createAuthor
//
// @Summary			Create new author
// @Description	Creates a new author
// @ID					authors-create-author
// @Tags				authors
// @Accept			json
// @Produce			json
// @Param				data						body 				AuthorToCreate			true	"Data for the author to create"
// @Success			201							string			string										"ID of the created author"
// @Failure			400							{object}		ErrorBadRequest						"Bad Request"
// @Failure			401							{object}		ErrorUnauthorized					"Unauthorized"
// @Failure 		500							{object}		ErrorInternalServerError	"Internal Server Error"
// @Router			/authors				[post]
func (server *Server) createAuthor(ctx *gin.Context) {
	var authorBody db.AuthorToCreate
	if err := ctx.ShouldBindJSON(&authorBody); err != nil {
		fmt.Println(err)
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	authorID, err := server.store.CreateAuthor(ctx, authorBody)
	if err != nil {
		if strings.HasPrefix(err.Error(), "author with name") { // TODO: Find better way to check error types (enum?)
			NewErrorBadRequest(err).Send(ctx)
			return
		}

		NewErrorInternalServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, authorID)
}

// getAuthorByID
//
// @Summary			Get one author by ID
// @Description	One author, which matches the ID, is returned
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
		return
	}

	ctx.JSON(http.StatusOK, author)
}

// patchAuthorByID
//
// @Summary			Patch one author by ID
// @Description	One author, which matches the ID, is modified with the provided patch
// @ID					authors-patch-author-by-id
// @Tags				authors
// @Accept			json
// @Produce			json
// @Param				id							path 				int									true	"ID of the desired author to patch"
// @Param				data						body 				AuthorUpdate				true	"Patch for modifying the author"
// @Success			200
// @Failure			400							{object}		ErrorBadRequest						"Bad Request"
// @Failure			401							{object}		ErrorUnauthorized					"Unauthorized"
// @Failure			404							{object}		ErrorNotFound							"Not Found"
// @Router			/authors/{id}		[patch]
func (server *Server) patchAuthorByID(ctx *gin.Context) {
	var uriParam getByIDRequest
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	var authorPatch db.AuthorUpdate
	if err := ctx.ShouldBindJSON(&authorPatch); err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	modifiedCount, err := server.store.UpdateAuthorByID(ctx, uriParam.ID, authorPatch)
	if err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	if modifiedCount < 1 {
		NewErrorNotFound(err).Send(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

// deleteAuthorByID
//
// @Summary			Delete one author by ID
// @Description	One author, which matches the ID, is deleted
// @ID					authors-delete-author-by-id
// @Tags				authors
// @Accept			json
// @Produce			json
// @Param				id							path 				int									true	"ID of the desired author to patch"
// @Success			200
// @Failure			400							{object}		ErrorBadRequest						"Bad Request"
// @Failure			401							{object}		ErrorUnauthorized					"Unauthorized"
// @Failure			404							{object}		ErrorNotFound							"Not Found"
// @Router			/authors/{id}		[delete]
func (server *Server) deleteAuthorByID(ctx *gin.Context) {
	var uriParam getByIDRequest
	if err := ctx.ShouldBindUri(&uriParam); err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	deleteCount, err := server.store.DeleteAuthorByID(ctx, uriParam.ID)
	if err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	if deleteCount < 1 {
		NewErrorNotFound(err).Send(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
