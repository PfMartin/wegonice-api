package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

const createFormFileName string = "image"

var allowedFileTypes = []string{".png", ".jpg"}

// saveImage
//
// @Summary			Saves an image
// @Description	Saves an image to the filesystem. Send a request with
// @Description `const formData = new FormData();`
// @Description `formData.append('image', image);`
// @Description Add the header `'ContentType': 'multipart/form-data';`
// @ID					images-save
// @Tags				images
// @Accept			multipart/form-data
// @Produce			json
// @Param				authorization	header			string							false	"Authorization header for bearer token"
// @Param				page_size			query 			int									true	"Number of elements in one page"
// @Success			200
// @Failure			400						{object}		ErrorBadRequest						"Bad Request"
// @Failure			401						{object}		ErrorUnauthorized					"Unauthorized"
// @Failure 		500						{object}		ErrorInternalServerError	"Internal Server Error"
// @Router			/images				[post]
func (s *Server) SaveImage(ctx *gin.Context) {
	file, err := ctx.FormFile(createFormFileName)
	if err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	extension := filepath.Ext(file.Filename)
	if !slices.Contains(allowedFileTypes, strings.ToLower(extension)) {
		fmt.Println(extension)
		NewErrorBadRequest(fmt.Errorf("file type %s is not supported", extension)).Send(ctx)
		return
	}

	destination := fmt.Sprintf("%s/%s", s.config.imagesDepotPath, file.Filename)
	if _, err := os.Stat(destination); err == nil {
		NewErrorBadRequest(fmt.Errorf("file already exists")).Send(ctx)
		return
	}

	if err = ctx.SaveUploadedFile(file, destination); err != nil {
		NewErrorInternalServerError(err).Send(ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

// getImage
//
// @Summary			Gets an image
// @Description	Gets an image with the given image name from the file system
// @ID					images-get
// @Tags				images
// @Accept			json
// @Produce			json
// @Param				authorization				header			string							false	"Authorization header for bearer token"
// @Param				page_size						query 			int									true	"Number of elements in one page"
// @Success			200									{array}			byte											"Image for the given image name"
// @Failure			400									{object}		ErrorBadRequest						"Bad Request"
// @Failure			401									{object}		ErrorUnauthorized					"Unauthorized"
// @Failure 		500									{object}		ErrorInternalServerError	"Internal Server Error"
// @Router			/images/{imagename}	[get]
func (s *Server) GetImage(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

// deleteImage
//
// @Summary			Deletes an image
// @Description	Deletes an image with the given image name from the file system
// @ID					images-delete
// @Tags				images
// @Accept			json
// @Produce			json
// @Param				authorization				header			string							false	"Authorization header for bearer token"
// @Param				page_size						query 			int									true	"Number of elements in one page"
// @Success			200									{array}			byte											"Image for the given image name"
// @Failure			400									{object}		ErrorBadRequest						"Bad Request"
// @Failure			401									{object}		ErrorUnauthorized					"Unauthorized"
// @Failure 		500									{object}		ErrorInternalServerError	"Internal Server Error"
// @Router			/images/{imagename}	[delete]
func (s *Server) DeleteImage(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
