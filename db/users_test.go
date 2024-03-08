package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) {
	conf := getDatabaseConfiguration(t)

	dbClient, _ := NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	require.NotNil(t, dbClient)

	user := User{
		Email:      "test@test.com",
		Password:   "testpassword",
		Role:       "User",
		CreatedAt:  time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	}

	handler := NewUsersHandler(dbClient, conf.DBName)

	ctx := context.Background()
	insertedId, err := handler.CreateUser(ctx, user)

	require.NoError(t, err)
	require.False(t, insertedId.IsZero())
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
