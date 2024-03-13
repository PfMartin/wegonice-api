package db

import (
	"context"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/stretchr/testify/require"
)

func getAuthorCollection(t *testing.T) *AuthorCollection {
	t.Helper()

	conf := getDatabaseConfiguration(t)

	dbClient, _ := NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	require.NotNil(t, dbClient)

	coll := NewAuthorCollection(dbClient, conf.DBName)

	return coll
}

func createRandomAuthor(t *testing.T, authorColl *AuthorCollection, userID string) Author {
	t.Helper()

	author := Author{
		FirstName:    util.RandomString(6),
		LastName:     util.RandomString(6),
		Name:         util.RandomString(6),
		WebsiteURL:   util.RandomString(10),
		InstagramURL: util.RandomString(10),
		YoutubeURL:   util.RandomString(10),
		ImageName:    util.RandomString(10),
		UserID:       userID,
	}

	insertedAuthorID, err := authorColl.CreateAuthor(context.Background(), author)
	require.NoError(t, err)
	require.False(t, insertedAuthorID.IsZero())

	authorID := insertedAuthorID.Hex()

	return Author{
		ID:           authorID,
		FirstName:    author.FirstName,
		LastName:     author.LastName,
		Name:         author.Name,
		WebsiteURL:   author.WebsiteURL,
		InstagramURL: author.InstagramURL,
		YoutubeURL:   author.YoutubeURL,
		ImageName:    author.ImageName,
		UserID:       author.UserID,
		CreatedAt:    time.Now().Unix(),
		ModifiedAt:   time.Now().Unix(),
	}
}

func TestUnitCreateAuthor(t *testing.T) {
	user := createRandomUser(t, getUserCollection(t))
	authorColl := getAuthorCollection(t)

	t.Run("Creates a new author and throws an error when the same author should be created again", func(t *testing.T) {
		author := createRandomAuthor(t, authorColl, user.ID)

		_, err := authorColl.CreateAuthor(context.Background(), author)
		require.Error(t, err)
	})
}

func TestUnitGetAllAuthors(t *testing.T) {
	user := createRandomUser(t, getUserCollection(t))
	authorColl := getAuthorCollection(t)

	for i := 0; i < 10; i++ {
		_ = createRandomAuthor(t, authorColl, user.ID)
	}

	pagination := Pagination{
		PageID:   1,
		PageSize: 5,
	}

	t.Run("Gets all authors with pagination", func(t *testing.T) {
		authors, err := authorColl.GetAllAuthors(context.Background(), pagination)
		require.NoError(t, err)
		require.NotEmpty(t, authors)

		require.Equal(t, int(pagination.PageSize), len(authors))
	})
}

func TestUnitGetAuthorByID(t *testing.T) {
	user := createRandomUser(t, getUserCollection(t))
	authorColl := getAuthorCollection(t)
	createdAuthor := createRandomAuthor(t, authorColl, user.ID)

	testCases := []struct {
		name           string
		authorID       string
		hasError       bool
		expectedAuthor Author
	}{
		{
			name:           "Success",
			authorID:       createdAuthor.ID,
			hasError:       false,
			expectedAuthor: createdAuthor,
		},
		{
			name:     "Fail with invalid authorID",
			authorID: "test",
			hasError: true,
		},
		{
			name:     "Fail with authorID not found",
			authorID: "659c00751f717854f690270d",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotAuthor, err := authorColl.GetAuthorByID(context.Background(), tc.authorID)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedAuthor, gotAuthor)
		})
	}
}

func TestUnitUpdateAuthorByID(t *testing.T) {
	user := createRandomUser(t, getUserCollection(t))
	authorColl := getAuthorCollection(t)
	createdAuthor := createRandomAuthor(t, authorColl, user.ID)

	authorUpdate := Author{
		Name:         util.RandomString(4),
		LastName:     util.RandomString(6),
		FirstName:    util.RandomString(6),
		WebsiteURL:   util.RandomString(6),
		InstagramURL: util.RandomString(6),
		YoutubeURL:   util.RandomString(6),
		ImageName:    util.RandomString(10),
	}

	testCases := []struct {
		name          string
		authorID      string
		authorUpdate  Author
		hasError      bool
		modifiedCount int64
	}{
		{
			name:          "Success",
			authorID:      createdAuthor.ID,
			authorUpdate:  authorUpdate,
			hasError:      false,
			modifiedCount: 1,
		},
		{
			name:          "Fail with invalid authorID",
			authorID:      "test",
			hasError:      true,
			modifiedCount: 0,
		},
		{
			name:          "Fail with authorID not found",
			authorID:      "659c00751f717854f690270d",
			hasError:      false,
			modifiedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			modifiedCount, err := authorColl.UpdateAuthorByID(context.Background(), tc.authorID, tc.authorUpdate)
			require.Equal(t, tc.modifiedCount, modifiedCount)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if modifiedCount < 1 {
				return
			}

			updatedAuthor, err := authorColl.GetAuthorByID(context.Background(), tc.authorID)
			require.NoError(t, err)

			expectedAuthor := Author{
				ID:           createdAuthor.ID,
				Name:         authorUpdate.Name,
				LastName:     authorUpdate.LastName,
				FirstName:    authorUpdate.FirstName,
				WebsiteURL:   authorUpdate.WebsiteURL,
				InstagramURL: authorUpdate.InstagramURL,
				YoutubeURL:   authorUpdate.YoutubeURL,
				ImageName:    authorUpdate.ImageName,
				CreatedAt:    createdAuthor.CreatedAt,
				ModifiedAt:   time.Now().Unix(),
			}

			require.Equal(t, expectedAuthor.ID, updatedAuthor.ID)
			require.Equal(t, expectedAuthor.Name, updatedAuthor.Name)
			require.Equal(t, expectedAuthor.LastName, updatedAuthor.LastName)
			require.Equal(t, expectedAuthor.FirstName, updatedAuthor.FirstName)
			require.Equal(t, expectedAuthor.WebsiteURL, updatedAuthor.WebsiteURL)
			require.Equal(t, expectedAuthor.InstagramURL, updatedAuthor.InstagramURL)
			require.Equal(t, expectedAuthor.YoutubeURL, updatedAuthor.YoutubeURL)
			require.Equal(t, expectedAuthor.ImageName, updatedAuthor.ImageName)
			require.WithinDuration(t, time.Unix(expectedAuthor.CreatedAt, 0), time.Unix(updatedAuthor.CreatedAt, 0), 1*time.Second)
			require.WithinDuration(t, time.Unix(expectedAuthor.ModifiedAt, 0), time.Unix(updatedAuthor.ModifiedAt, 0), 1*time.Second)
		})
	}
}

func TestDeleteAuthorByID(t *testing.T) {
	user := createRandomUser(t, getUserCollection(t))
	authorColl := getAuthorCollection(t)
	createdAuthor := createRandomAuthor(t, authorColl, user.ID)

	testCases := []struct {
		name        string
		authorID    string
		hasError    bool
		deleteCount int64
	}{
		{
			name:        "Success",
			authorID:    createdAuthor.ID,
			hasError:    false,
			deleteCount: 1,
		},
		{
			name:        "Fail with invalid authorID",
			authorID:    "test",
			hasError:    true,
			deleteCount: 0,
		},
		{
			name:        "Fail with authorID not found",
			authorID:    "659c00751f717854f690270d",
			hasError:    false,
			deleteCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deleteCount, err := authorColl.DeleteAuthorByID(context.Background(), tc.authorID)
			require.Equal(t, tc.deleteCount, deleteCount)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tc.deleteCount, deleteCount)
		})
	}
}
