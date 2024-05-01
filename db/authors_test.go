package db

import (
	"context"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/stretchr/testify/require"
)

func getMongoDBStore(t *testing.T) *MongoDBStore {
	t.Helper()

	conf := getDatabaseConfiguration(t)

	dbClient, _ := NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	store := NewMongoDBStore(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	require.NotNil(t, dbClient)

	return store
}

func createRandomAuthor(t *testing.T, store *MongoDBStore, userID string) Author {
	t.Helper()

	author := AuthorToCreate{
		FirstName:    util.RandomString(6),
		LastName:     util.RandomString(6),
		Name:         util.RandomString(6),
		WebsiteURL:   util.RandomString(10),
		InstagramURL: util.RandomString(10),
		YoutubeURL:   util.RandomString(10),
		ImageName:    util.RandomString(10),
		UserID:       userID,
	}

	insertedAuthorID, err := store.CreateAuthor(context.Background(), author)
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
	store := getMongoDBStore(t)
	user := createRandomUser(t, store)

	t.Run("Creates a new author and throws an error when the same author should be created again", func(t *testing.T) {
		author := createRandomAuthor(t, store, user.ID)

		authorToCreate := AuthorToCreate{
			FirstName:    author.FirstName,
			LastName:     author.LastName,
			Name:         author.Name,
			WebsiteURL:   author.WebsiteURL,
			InstagramURL: author.InstagramURL,
			YoutubeURL:   author.YoutubeURL,
			ImageName:    author.ImageName,
			UserID:       author.UserID,
		}
		_, err := store.CreateAuthor(context.Background(), authorToCreate)
		require.Error(t, err)
	})
}

func TestUnitGetAllAuthors(t *testing.T) {
	store := getMongoDBStore(t)
	user := createRandomUser(t, store)

	for i := 0; i < 10; i++ {
		_ = createRandomAuthor(t, store, user.ID)
	}

	pagination := Pagination{
		PageID:   1,
		PageSize: 5,
	}

	t.Run("Gets all authors with pagination", func(t *testing.T) {
		ctx := context.Background()
		authors, err := store.GetAllAuthors(ctx, pagination)
		require.NoError(t, err)
		require.NotEmpty(t, authors)

		require.Equal(t, int(pagination.PageSize), len(authors))

		for _, author := range authors {
			require.NotEmpty(t, author.UserCreated.Email)
			require.NotEmpty(t, author.UserCreated.ID)
			require.Empty(t, author.UserID)
		}
	})
}

func TestUnitGetAuthorByID(t *testing.T) {
	store := getMongoDBStore(t)
	user := createRandomUser(t, store)

	createdAuthor := createRandomAuthor(t, store, user.ID)

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
			authorID: "659c00751f7178dff690270d",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotAuthor, err := store.GetAuthorByID(context.Background(), tc.authorID)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			tc.expectedAuthor.UserCreated = User{
				ID:    user.ID,
				Email: user.Email,
			}

			tc.expectedAuthor.UserID = ""

			require.NoError(t, err)
			require.Equal(t, tc.expectedAuthor, gotAuthor)
		})
	}
}

func TestUnitUpdateAuthorByID(t *testing.T) {
	store := getMongoDBStore(t)
	user := createRandomUser(t, store)

	createdAuthor := createRandomAuthor(t, store, user.ID)

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
			modifiedCount, err := store.UpdateAuthorByID(context.Background(), tc.authorID, tc.authorUpdate)
			require.Equal(t, tc.modifiedCount, modifiedCount)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if modifiedCount < 1 {
				return
			}

			updatedAuthor, err := store.GetAuthorByID(context.Background(), tc.authorID)
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
			require.WithinDuration(t, time.Unix(expectedAuthor.CreatedAt, 0), time.Unix(updatedAuthor.CreatedAt, 0), 5*time.Second)
			require.WithinDuration(t, time.Unix(expectedAuthor.ModifiedAt, 0), time.Unix(updatedAuthor.ModifiedAt, 0), 5*time.Second)
		})
	}
}

func TestUnitDeleteAuthorByID(t *testing.T) {
	store := getMongoDBStore(t)

	user := createRandomUser(t, store)
	createdAuthor := createRandomAuthor(t, store, user.ID)
	recipeUser := createRandomUser(t, store)
	recipeAuthor := createRandomAuthor(t, store, recipeUser.ID)

	testCases := []struct {
		name                 string
		authorID             string
		hasError             bool
		hasReferencingRecipe bool
		deleteCount          int64
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
			authorID:    "659c00751f717df4f690270d",
			hasError:    false,
			deleteCount: 0,
		},
		{
			name:                 "Fail with author referenced in at least one recipes",
			authorID:             recipeAuthor.ID,
			hasError:             true,
			hasReferencingRecipe: true,
			deleteCount:          0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.hasReferencingRecipe {
				createRandomRecipe(t, store, recipeUser.ID, recipeAuthor.ID)
			}

			deleteCount, err := store.DeleteAuthorByID(context.Background(), tc.authorID)
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
