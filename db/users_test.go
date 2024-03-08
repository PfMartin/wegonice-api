package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	conf := getDatabaseConfiguration(t)

	dbClient, cancel := NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	defer cancel()
	require.NotNil(t, dbClient)

	user := User{
		Email:      "test@test.com",
		Password:   "testpassword",
		Role:       "User",
		CreatedAt:  time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	}

	handler := NewUsersHandler(dbClient, conf.DBName)

	ctx, cancel := context.WithCancel(context.Background())
	insertedId, err := handler.CreateUser(ctx, user)

	require.NoError(t, err)
	require.False(t, insertedId.IsZero())
	cancel()
}
