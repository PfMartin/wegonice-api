package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PfMartin/wegonice-api/db"
	mock_db "github.com/PfMartin/wegonice-api/db/mock"
	"github.com/PfMartin/wegonice-api/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func randomAuthor(t *testing.T) db.Author {
	t.Helper()

	userID := primitive.NewObjectID().Hex()

	return db.Author{
		ID:           primitive.NewObjectID().Hex(),
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
	}
}

func TestUnitListAuthors(t *testing.T) {
	var authors []db.Author
	for range 10 {
		authors = append(authors, randomAuthor(t))
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

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUnitGetAuthorByID(t *testing.T) {
	var authors []db.Author
	for range 2 {
		authors = append(authors, randomAuthor(t))
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
