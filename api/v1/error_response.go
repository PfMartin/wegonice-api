package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorBadRequest struct {
	StatusText string `json:"statusText" example:"Bad Request"`
	StatusCode int    `json:"statusCode" example:"400"`
	Message    string `json:"message" example:"Failed to parse data"`
} // @name ErrorBadRequest

func NewErrorBadRequest(err error) *ErrorBadRequest {
	return &ErrorBadRequest{
		StatusText: http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Message:    err.Error(),
	}
}

func (err *ErrorBadRequest) Send(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(err.StatusCode, err)
}

type ErrorInternalServerError struct {
	StatusText string `json:"statusText" example:"Internal Server Error"`
	StatusCode int    `json:"statusCode" example:"500"`
	Message    string `json:"message" example:"An internal server error occurred"`
} // @name ErrorInternalServerError

func NewErrorInternalServerError(err error) *ErrorInternalServerError {
	return &ErrorInternalServerError{
		StatusText: http.StatusText(http.StatusInternalServerError),
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
	}
}

func (err *ErrorInternalServerError) Send(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(err.StatusCode, err)
}

type ErrorNotAcceptable struct {
	StatusText string `json:"statusText" example:"Not Acceptable"`
	StatusCode int    `json:"statusCode" example:"406"`
	Message    string `json:"message" example:"Provided input is not acceptable"`
} // @name ErrorNotAcceptable

func NewErrorNotAcceptable(err error) *ErrorNotAcceptable {
	return &ErrorNotAcceptable{
		StatusText: http.StatusText(http.StatusNotAcceptable),
		StatusCode: http.StatusNotAcceptable,
		Message:    err.Error(),
	}
}

func (err *ErrorNotAcceptable) Send(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(err.StatusCode, err)
}

type ErrorNotFound struct {
	StatusText string `json:"statusText" example:"Not Found"`
	StatusCode int    `json:"statusCode" example:"404"`
	Message    string `json:"message" example:"Could not find requested data"`
} // @name ErrorNotFound

func NewErrorNotFound(err error) *ErrorNotFound {
	return &ErrorNotFound{
		StatusText: http.StatusText(http.StatusNotFound),
		StatusCode: http.StatusNotFound,
		Message:    err.Error(),
	}
}

func (err *ErrorNotFound) Send(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(err.StatusCode, err)
}

type ErrorUnauthorized struct {
	StatusText string `json:"statusText" example:"Unauthorized"`
	StatusCode int    `json:"statusCode" example:"401"`
	Message    string `json:"message" example:"Unauthorized for retrieving this information"`
} // @name ErrorUnauthorized

func NewErrorUnauthorized(err error) *ErrorUnauthorized {
	return &ErrorUnauthorized{
		StatusText: http.StatusText(http.StatusUnauthorized),
		StatusCode: http.StatusUnauthorized,
		Message:    err.Error(),
	}
}

func (err *ErrorUnauthorized) Send(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(err.StatusCode, err)
}
