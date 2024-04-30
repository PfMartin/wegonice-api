package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type ErrorResponse struct {
	StatusText string `json:"statusText"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func TestUnitErrorResponses(t *testing.T) {
	router := gin.Default()

	router.GET("/bad_request", func(ctx *gin.Context) {
		NewErrorBadRequest(fmt.Errorf("bad request")).Send(ctx)
	})
	router.GET("/internal_server_error", func(ctx *gin.Context) {
		NewErrorInternalServerError(fmt.Errorf("internal server error")).Send(ctx)
	})
	router.GET("/not_acceptable", func(ctx *gin.Context) {
		NewErrorNotAcceptable(fmt.Errorf("not acceptable")).Send(ctx)
	})
	router.GET("/not_found", func(ctx *gin.Context) {
		NewErrorNotFound(fmt.Errorf("not found")).Send(ctx)
	})
	router.GET("/unauthorized", func(ctx *gin.Context) {
		NewErrorUnauthorized(fmt.Errorf("unauthorized")).Send(ctx)
	})

	testCases := []struct {
		name               string
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name:               "bad_request",
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "bad request",
		},
		{
			name:               "internal_server_error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedMessage:    "internal server error",
		},
		{
			name:               "not_acceptable",
			expectedStatusCode: http.StatusNotAcceptable,
			expectedMessage:    "not acceptable",
		},
		{
			name:               "not_found",
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "not found",
		},
		{
			name:               "unauthorized",
			expectedStatusCode: http.StatusUnauthorized,
			expectedMessage:    "unauthorized",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			url := "/" + tc.name
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			router.ServeHTTP(recorder, request)
			require.Equal(t, tc.expectedStatusCode, recorder.Code)

			var errorBody ErrorResponse
			err = json.NewDecoder(recorder.Body).Decode(&errorBody)
			require.NoError(t, err)

			require.Equal(t, tc.expectedMessage, errorBody.Message)
			require.Equal(t, http.StatusText(tc.expectedStatusCode), errorBody.StatusText)
		})
	}
}
