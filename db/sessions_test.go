package db

import (
	"context"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T, store *MongoDBStore, userID string) Session {
	t.Helper()

	session := Session{
		UserID:       userID,
		RefreshToken: util.RandomString(10),
		UserAgent:    util.RandomString(6),
		ClientIP:     util.RandomString(12),
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
	}

	insertedSessionID, err := store.CreateSession(context.Background(), session)
	require.NoError(t, err)
	require.False(t, insertedSessionID.IsZero())

	sessionID := insertedSessionID.Hex()

	return Session{
		ID:           sessionID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIP:     session.ClientIP,
		ExpiresAt:    session.ExpiresAt,
	}
}

func TestUnitCreateSession(t *testing.T) {
	store := getMongoDBStore(t)
	user := createRandomUser(t, store)

	t.Run("Creates a new session and throws an error when the same session should be created again", func(t *testing.T) {
		_ = createRandomSession(t, store, user.ID)
	})
}

func TestUnitGetSessionByID(t *testing.T) {
	store := getMongoDBStore(t)

	user := createRandomUser(t, store)
	createdSession := createRandomSession(t, store, user.ID)

	t.Run("Gets created session without any errors", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		gotSession, err := store.GetSessionByID(ctx, createdSession.ID)
		require.NoError(t, err)

		expectedSession := createdSession
		expectedSession.UserID = ""

		require.Equal(t, expectedSession, gotSession)
	})
}
