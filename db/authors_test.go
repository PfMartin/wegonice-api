package db

import (
	"context"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/stretchr/testify/require"
)

func getAuthorCollection(t *testing.T) *AuthorCollection {
	conf := getDatabaseConfiguration(t)

	dbClient, _ := NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	require.NotNil(t, dbClient)

	coll := NewAuthorCollection(dbClient, conf.DBName)

	return coll
}

func createRandomAuthor(t *testing.T, authorColl *AuthorCollection, userID string) Author {
	author := Author{
		FirstName:    util.RandomString(6),
		LastName:     util.RandomString(6),
		Name:         util.RandomString(6),
		WebsiteURL:   util.RandomString(10),
		InstagramURL: util.RandomString(10),
		YouTubeURL:   util.RandomString(10),
		ImageBase64:  util.RandomString(10),
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
		YouTubeURL:   author.YouTubeURL,
		ImageBase64:  author.ImageBase64,
		UserID:       author.UserID,
		CreatedAt:    time.Now().UnixMilli(),
		ModifiedAt:   time.Now().UnixMilli(),
	}
}

func TestCreateAuthor(t *testing.T) {
	user := createRandomUser(t, getUserCollection(t))
	authorColl := getAuthorCollection(t)

	t.Run("Creates a new author and throws an error when the same author should be created again", func(t *testing.T) {
		author := createRandomAuthor(t, authorColl, user.ID)

		_, err := authorColl.CreateAuthor(context.Background(), author)
		require.Error(t, err)
	})
}

func TestGetAllAuthors(t *testing.T) {
	user := createRandomUser(t, getUserCollection(t))
	authorColl := getAuthorCollection(t)

	for range 7 {
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
