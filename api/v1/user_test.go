package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PfMartin/wegonice-api/db"
	mock_db "github.com/PfMartin/wegonice-api/db/mock"
	"github.com/PfMartin/wegonice-api/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (db.User, string) {
	t.Helper()

	password := util.RandomString(8)

	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		Email:        util.RandomEmail(),
		PasswordHash: hashedPassword,
	}

	return user, password
}

func TestUnitRegisterUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockDBStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/api/v1/auth/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
