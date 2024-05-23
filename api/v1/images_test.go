package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	mock_db "github.com/PfMartin/wegonice-api/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUnitSaveImage(t *testing.T) {
	testCases := []struct {
		name         string
		filename     string
		responseCode int
	}{
		{
			name:         "Success with .png file",
			filename:     "test_image.png",
			responseCode: http.StatusOK,
		},
		{
			name:         "Success with .jpg file",
			filename:     "test_image.jpg",
			responseCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			w := multipart.NewWriter(&buf)

			testImagePath := fmt.Sprintf("../../images/test_images/%s", tc.filename)

			file, err := os.Open(testImagePath)
			require.NoError(t, err)
			defer file.Close()

			fw, err := w.CreateFormFile(createFormFileName, tc.filename)
			require.NoError(t, err)

			_, err = io.Copy(fw, file)
			require.NoError(t, err)

			w.Close()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockDBStore(ctrl)

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			url := "/api/v1/images"
			request, err := http.NewRequest(http.MethodPost, url, &buf)
			require.NoError(t, err)
			request.Header.Set("Content-Type", w.FormDataContentType())

			user, _ := randomUser(t)
			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, tc.responseCode, recorder.Code)
		})
	}
}

func TestUnitServeImage(t *testing.T) {
	testCases := []struct {
		name         string
		imageName    string
		responseCode int
	}{
		{
			name:         "Success with .png image",
			imageName:    "test_image.png",
			responseCode: http.StatusOK,
		},
		{
			name:         "Success with .jpg image",
			imageName:    "test_image.jpg",
			responseCode: http.StatusOK,
		},
		{
			name:         "Fail with non-existing image",
			imageName:    "non-existing.jpg",
			responseCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockDBStore(ctrl)

			server := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/images/%s", tc.imageName)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			user, _ := randomUser(t)
			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			require.Equal(t, tc.responseCode, recorder.Code)
		})
	}
}
