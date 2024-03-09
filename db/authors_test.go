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

func createRandomAuthor(t *testing.T, authorColl *AuthorCollection) Author {
	userColl := getUserCollection(t)

	user := createRandomUser(t, userColl)

	author := Author{
		FirstName:    util.RandomString(6),
		LastName:     util.RandomString(6),
		Name:         util.RandomString(6),
		WebsiteURL:   util.RandomString(10),
		InstagramURL: util.RandomString(10),
		YouTubeURL:   util.RandomString(10),
		ImageBase64:  util.RandomString(10),
		UserID:       user.ID,
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
	authorColl := getAuthorCollection(t)
	createRandomAuthor(t, authorColl)
}
