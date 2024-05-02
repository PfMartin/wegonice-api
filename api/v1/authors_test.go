package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/db"
	mock_db "github.com/PfMartin/wegonice-api/db/mock"
	"github.com/PfMartin/wegonice-api/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func randomAuthor(t *testing.T) (db.Author, primitive.ObjectID) {
	t.Helper()

	userID := primitive.NewObjectID().Hex()
	authorID := primitive.NewObjectID()

	return db.Author{
		ID:           authorID.Hex(),
		FirstName:    util.RandomString(6),
		LastName:     util.RandomString(6),
		Name:         util.RandomString(6),
		WebsiteURL:   util.RandomString(10),
		InstagramURL: util.RandomString(10),
		YoutubeURL:   util.RandomString(10),
		ImageName:    util.RandomString(10),
		RecipeCount:  int(util.RandomInt(0, 100)),
		UserID:       userID,
		UserCreated: db.User{
			ID:    userID,
			Email: util.RandomEmail(),
		},
	}, authorID
}

func TestUnitListAuthors(t *testing.T) {
	user, _ := randomUser(t)
	var authors []db.Author
	for i := 0; i < 10; i++ {
		author, _ := randomAuthor(t)

		authors = append(authors, author)
	}

	testCases := []struct {
		name          string
		query         string
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "Success with pagination from 1 to 10",
			query: "?page_id=1&page_size=10",
			buildStubs: func(store *mock_db.MockDBStore) {
				pagination := db.Pagination{
					PageID:   1,
					PageSize: 10,
				}

				store.EXPECT().GetAllAuthors(gomock.Any(), pagination).Times(1).Return(authors, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var gotAuthors []AuthorResponse
				err := json.NewDecoder(recorder.Body).Decode(&gotAuthors)
				require.NoError(t, err)

				require.Equal(t, 10, len(gotAuthors))

				for i, expectedAuthor := range authors {
					requireAuthorComparison(t, expectedAuthor, gotAuthors[i])
				}
			},
		},
		{
			name:  "Success with pagination from 5 to 10",
			query: "?page_id=5&page_size=10",
			buildStubs: func(store *mock_db.MockDBStore) {
				pagination := db.Pagination{
					PageID:   5,
					PageSize: 10,
				}

				store.EXPECT().GetAllAuthors(gomock.Any(), pagination).Times(1).Return(authors[4:], nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var gotAuthors []AuthorResponse
				err := json.NewDecoder(recorder.Body).Decode(&gotAuthors)
				require.NoError(t, err)

				require.Equal(t, 6, len(gotAuthors))

				for i, expectedAuthor := range authors[4:] {
					requireAuthorComparison(t, expectedAuthor, gotAuthors[i])
				}
			},
		},
		{
			name:  "Fail with missing page_size",
			query: "?page_id=5",
			buildStubs: func(store *mock_db.MockDBStore) {
				pagination := db.Pagination{
					PageID: 5,
				}

				store.EXPECT().GetAllAuthors(gomock.Any(), pagination).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Fail with missing page_id",
			query: "?page_size=10",
			buildStubs: func(store *mock_db.MockDBStore) {
				pagination := db.Pagination{
					PageSize: 10,
				}

				store.EXPECT().GetAllAuthors(gomock.Any(), pagination).Times(0)
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

			url := fmt.Sprintf("/api/v1/authors%s", tc.query)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUnitGetAuthorByID(t *testing.T) {
	user, _ := randomUser(t)
	var authors []db.Author
	for i := 0; i < 2; i++ {
		author, _ := randomAuthor(t)
		authors = append(authors, author)
	}

	testCases := []struct {
		name          string
		id            string
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success getting the second author",
			id:   authors[1].ID,
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().GetAuthorByID(gomock.Any(), authors[1].ID).Times(1).Return(authors[1], nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var gotAuthor AuthorResponse
				err := json.NewDecoder(recorder.Body).Decode(&gotAuthor)
				require.NoError(t, err)

				requireAuthorComparison(t, authors[1], gotAuthor)
			},
		},
		{
			name: "Fail with non-existent ID",
			id:   "notexisting",
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().GetAuthorByID(gomock.Any(), "notexisting").Times(1).Return(db.Author{}, fmt.Errorf("failed to find author with authorID: notexisting"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Fail with non-parsable ID",
			id:   "notexisting",
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().GetAuthorByID(gomock.Any(), "notexisting").Times(1).Return(db.Author{}, fmt.Errorf("failed to parse authorID"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			url := fmt.Sprintf("/api/v1/authors/%s", tc.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUnitCreateAuthor(t *testing.T) {
	user, _ := randomUser(t)
	author, primitiveID := randomAuthor(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success creating a new author",
			body: gin.H{
				"name":      author.Name,
				"firstName": author.FirstName,
				"lastName":  author.LastName,
				"userId":    author.UserID,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateAuthor(gomock.Any(), db.AuthorToCreate{
					Name:      author.Name,
					FirstName: author.FirstName,
					LastName:  author.LastName,
					UserID:    author.UserID,
				}).Times(1).Return(primitiveID, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Fail due to missing name",
			body: gin.H{
				"firstName": author.FirstName,
				"lastName":  author.LastName,
				"userId":    author.UserID,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateAuthor(gomock.Any(), db.AuthorToCreate{
					Name:      author.Name,
					FirstName: author.FirstName,
					LastName:  author.LastName,
					UserID:    author.UserID,
				}).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fail due to missing userID",
			body: gin.H{
				"name":      author.Name,
				"firstName": author.FirstName,
				"lastName":  author.LastName,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateAuthor(gomock.Any(), db.AuthorToCreate{
					Name:      author.Name,
					FirstName: author.FirstName,
					LastName:  author.LastName,
					UserID:    author.UserID,
				}).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fail due to already existing author",
			body: gin.H{
				"name":      author.Name,
				"firstName": author.FirstName,
				"lastName":  author.LastName,
				"userId":    author.UserID,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().CreateAuthor(gomock.Any(), db.AuthorToCreate{
					Name:      author.Name,
					FirstName: author.FirstName,
					LastName:  author.LastName,
					UserID:    author.UserID,
				}).Times(1).Return(primitive.NilObjectID, fmt.Errorf("author with name %s already exists", author.Name))
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

			url := "/api/v1/authors/"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUnitPatchAuthorByID(t *testing.T) {
	user, _ := randomUser(t)
	author, _ := randomAuthor(t)
	nonMatchingID := primitive.NewObjectID().Hex()

	fullAuthorPatch := db.AuthorUpdate{
		Name:         "New author name",
		FirstName:    "New first name",
		LastName:     "New last name",
		WebsiteURL:   "new website url",
		InstagramURL: "new instagram url",
		YoutubeURL:   "new youtube url",
		ImageName:    "new image name",
	}

	testCases := []struct {
		name          string
		id            string
		body          gin.H
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success with full update of the author",
			id:   author.ID,
			body: gin.H{
				"name":         fullAuthorPatch.Name,
				"firstName":    fullAuthorPatch.FirstName,
				"lastName":     fullAuthorPatch.LastName,
				"websiteUrl":   fullAuthorPatch.WebsiteURL,
				"instagramUrl": fullAuthorPatch.InstagramURL,
				"youtubeUrl":   fullAuthorPatch.YoutubeURL,
				"imageName":    fullAuthorPatch.ImageName,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateAuthorByID(gomock.Any(), author.ID, fullAuthorPatch).Times(1).Return(int64(1), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Success with partial update of the author",
			id:   author.ID,
			body: gin.H{
				"name":      fullAuthorPatch.Name,
				"firstName": fullAuthorPatch.FirstName,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateAuthorByID(gomock.Any(), author.ID, db.AuthorUpdate{
					Name:      fullAuthorPatch.Name,
					FirstName: fullAuthorPatch.FirstName,
				}).Times(1).Return(int64(1), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Fail due to missing id",
			id:   "",
			body: gin.H{
				"name": fullAuthorPatch.Name,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateAuthorByID(gomock.Any(), author.ID, db.AuthorUpdate{
					Name: fullAuthorPatch.Name,
				}).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Fail due to missing body",
			id:   author.ID,
			body: gin.H{},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateAuthorByID(gomock.Any(), author.ID, db.AuthorUpdate{}).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fail due to the provided author not being valid",
			id:   "not-valid-id",
			body: gin.H{
				"name": fullAuthorPatch.Name,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateAuthorByID(gomock.Any(), "not-valid-id", db.AuthorUpdate{
					Name: fullAuthorPatch.Name,
				}).Times(1).Return(int64(0), fmt.Errorf("failed to parse authorID"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Fail due to no matching author for author ID",
			id:   nonMatchingID,
			body: gin.H{
				"name": fullAuthorPatch.Name,
			},
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().UpdateAuthorByID(gomock.Any(), nonMatchingID, db.AuthorUpdate{
					Name: fullAuthorPatch.Name,
				}).Times(1).Return(int64(0), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := fmt.Sprintf("/api/v1/authors/%s", tc.id)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(data))
			require.NoError(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUnitDeleteAuthorByID(t *testing.T) {
	user, _ := randomUser(t)
	author, _ := randomAuthor(t)

	testCases := []struct {
		name          string
		id            string
		buildStubs    func(store *mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success deleting the author",
			id:   author.ID,
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().DeleteAuthorByID(gomock.Any(), author.ID).Times(1).Return(int64(1), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Fail due to missing id",
			id:   "",
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().DeleteAuthorByID(gomock.Any(), "").Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Fail due author id not matching an authors in the database",
			id:   "not-existing",
			buildStubs: func(store *mock_db.MockDBStore) {
				store.EXPECT().DeleteAuthorByID(gomock.Any(), "not-existing").Times(1).Return(int64(0), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := fmt.Sprintf("/api/v1/authors/%s", tc.id)

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Email, time.Minute)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireAuthorComparison(t *testing.T, expectedAuthor db.Author, gotAuthor AuthorResponse) {
	require.Equal(t, expectedAuthor.ID, gotAuthor.ID)
	require.Equal(t, expectedAuthor.FirstName, gotAuthor.FirstName)
	require.Equal(t, expectedAuthor.LastName, gotAuthor.LastName)
	require.Equal(t, expectedAuthor.Name, gotAuthor.Name)
	require.Equal(t, expectedAuthor.RecipeCount, gotAuthor.RecipeCount)
	require.Equal(t, expectedAuthor.WebsiteURL, gotAuthor.WebsiteURL)
	require.Equal(t, expectedAuthor.InstagramURL, gotAuthor.InstagramURL)
	require.Equal(t, expectedAuthor.YoutubeURL, gotAuthor.YoutubeURL)
	require.Equal(t, expectedAuthor.ImageName, gotAuthor.ImageName)
	require.Equal(t, expectedAuthor.UserID, gotAuthor.UserID)
	require.Equal(t, expectedAuthor.UserCreated.ID, gotAuthor.UserCreated.ID)
	require.Equal(t, expectedAuthor.UserCreated.Email, gotAuthor.UserCreated.Email)
}
