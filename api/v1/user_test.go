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
	"go.mongodb.org/mongo-driver/bson/primitive"

	"fmt"
)

func randomUser(t *testing.T) (db.User, string) {
	t.Helper()

	password := util.RandomString(8)

	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user := db.User{
		ID:           primitive.NewObjectID().Hex(),
		Email:        util.RandomEmail(),
		PasswordHash: hashedPassword,
	}

	return user, password
}

func TestUnitRegisterUser(t *testing.T) {
	user, password := randomUser(t)
	userID, err := primitive.ObjectIDFromHex(user.ID)
	require.NoError(t, err)

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
				userToCreate := db.User{
					Email:    user.Email,
					Password: password,
					IsActive: false,
				}

				store.EXPECT().CreateUser(gomock.Any(), userToCreate).Times(1).Return(userID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "Fails due to user with email already exists",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(primitive.NewObjectID(), fmt.Errorf("Duplicate user error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotAcceptable, recorder.Code)
			},
		},
		{
			name: "Fails due to missing email",
			body: gin.H{
				"password": password,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fails due to password too short",
			body: gin.H{
				"email":    user.Email,
				"password": "short",
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fails due to missing password",
			body: gin.H{
				"email": user.Email,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

func TestUnitLoginUser(t *testing.T) {
	user, password := randomUser(t)
	// userID, err := primitive.ObjectIDFromHex(user.ID)
	// require.NoError(t, err)

	sessionID := primitive.NewObjectID()

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
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(1).Return(sessionID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
			},
		},
		{
			name: "Fails due to missing email address",
			body: gin.H{
				"password": password,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fails due to missing password",
			body: gin.H{
				"email": user.Email,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fails due to wrong password",
			body: gin.H{
				"email":    user.Email,
				"password": "wrongPassword",
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(user, nil)
				store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			url := "/api/v1/auth/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
