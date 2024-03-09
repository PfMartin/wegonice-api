package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/stretchr/testify/require"
)

func getUsersHandler(t *testing.T) *CollectionHandler {
	conf := getDatabaseConfiguration(t)

	dbClient, _ := NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	require.NotNil(t, dbClient)

	handler := NewUsersHandler(dbClient, conf.DBName)

	return handler
}

func createRandomUser(t *testing.T, handler *CollectionHandler) User {
	user := User{
		Email:      util.RandomEmail(),
		Password:   util.RandomString(6),
		Role:       UserRole,
		CreatedAt:  time.Now().UnixMilli(),
		ModifiedAt: time.Now().UnixMilli(),
	}

	ctx := context.Background()
	insertedID, err := handler.CreateUser(ctx, user)

	require.NoError(t, err)
	require.False(t, insertedID.IsZero())

	userID := insertedID.Hex()

	return User{
		ID:         userID,
		Email:      user.Email,
		Password:   user.Password,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt,
		ModifiedAt: user.ModifiedAt,
	}
}

func TestCreateUser(t *testing.T) {
	handler := getUsersHandler(t)

	t.Run("Creates a new user", func(t *testing.T) {
		_ = createRandomUser(t, handler)
	})
}

func TestGetAllUsers(t *testing.T) {
	handler := getUsersHandler(t)

	for i := 0; i < 7; i++ {
		_ = createRandomUser(t, handler)
	}

	pagination := Pagination{
		PageID:   1,
		PageSize: 5,
	}

	t.Run("Gets all users with pagination", func(t *testing.T) {
		ctx := context.Background()
		users, err := handler.GetAllUsers(ctx, pagination)
		require.NoError(t, err)
		require.NotEmpty(t, users)

		require.Equal(t, int(pagination.PageSize), len(users))
	})
}

func TestGetUserByID(t *testing.T) {
	handler := getUsersHandler(t)

	createdUser := createRandomUser(t, handler)

	testCases := []struct {
		name         string
		userID       string
		hasError     bool
		expectedUser User
	}{
		{
			name:         "Success",
			userID:       createdUser.ID,
			hasError:     false,
			expectedUser: createdUser,
		},
		{
			name:     "Fail with invalid userID",
			userID:   "test",
			hasError: true,
		},
		{
			name:     "Fail with userID not found",
			userID:   "659c00751f717854f690270d",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotUser, err := handler.GetUserByID(context.Background(), tc.userID)

			fmt.Println(gotUser)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedUser, gotUser)
		})
	}
}

func TestUpdateUserByID(t *testing.T) {
	handler := getUsersHandler(t)

	createdUser := createRandomUser(t, handler)

	userUpdate := User{
		Email:    util.RandomEmail(),
		Password: util.RandomString(6),
	}

	testCases := []struct {
		name          string
		userID        string
		userUpdate    User
		hasError      bool
		modifiedCount int64
	}{
		{
			name:          "Success",
			userID:        createdUser.ID,
			userUpdate:    userUpdate,
			hasError:      false,
			modifiedCount: 1,
		},
		{
			name:          "Fail with invalid userID",
			userID:        "test",
			hasError:      true,
			modifiedCount: 0,
		},
		{
			name:          "Fail with userID not found",
			userID:        "659c00751f717854f690270d",
			hasError:      false,
			modifiedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			modifiedCount, err := handler.UpdateUserByID(context.Background(), tc.userID, tc.userUpdate)
			require.Equal(t, tc.modifiedCount, modifiedCount)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if modifiedCount < 1 {
				return
			}

			updatedUser, err := handler.GetUserByID(context.Background(), tc.userID)
			require.NoError(t, err)

			expectedUser := User{
				ID:         createdUser.ID,
				Email:      userUpdate.Email,
				Password:   userUpdate.Password,
				Role:       createdUser.Role,
				CreatedAt:  createdUser.CreatedAt,
				ModifiedAt: time.Now().UnixMilli(),
			}

			require.Equal(t, expectedUser.ID, updatedUser.ID)
			require.Equal(t, expectedUser.Email, updatedUser.Email)
			require.Equal(t, expectedUser.Password, updatedUser.Password)
			require.Equal(t, expectedUser.Role, updatedUser.Role)
			require.Equal(t, expectedUser.CreatedAt, updatedUser.CreatedAt)
			require.True(t, expectedUser.ModifiedAt <= updatedUser.ModifiedAt+100 || expectedUser.ModifiedAt >= updatedUser.ModifiedAt-100)
		})
	}
}

func TestDeleteUserByID(t *testing.T) {
	handler := getUsersHandler(t)

	createdUser := createRandomUser(t, handler)

	testCases := []struct {
		name        string
		userID      string
		hasError    bool
		deleteCount int64
	}{
		{
			name:        "Success",
			userID:      createdUser.ID,
			hasError:    false,
			deleteCount: 1,
		},
		{
			name:        "Fail with invalid userID",
			userID:      "test",
			hasError:    true,
			deleteCount: 0,
		},
		{
			name:        "Fail with userID not found",
			userID:      "659c00751f717854f690270d",
			hasError:    false,
			deleteCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deleteCount, err := handler.DeleteUserByID(context.Background(), tc.userID)
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
