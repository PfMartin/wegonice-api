package db

import (
	"context"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/stretchr/testify/require"
)

func getUserCollection(t *testing.T) *UserCollection {
	conf := getDatabaseConfiguration(t)

	dbClient, _ := NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	require.NotNil(t, dbClient)

	coll := NewUserCollection(dbClient, conf.DBName)

	return coll
}

func createRandomUser(t *testing.T, coll *UserCollection) User {
	user := User{
		Email:      util.RandomEmail(),
		Password:   util.RandomString(6),
		Role:       UserRole,
		CreatedAt:  time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	}

	ctx := context.Background()
	insertedID, err := coll.CreateUser(ctx, user)

	require.NoError(t, err)
	require.False(t, insertedID.IsZero())

	userID := insertedID.Hex()

	hashedPassword, err := util.HashPassword(user.Password)
	require.NoError(t, err)

	return User{
		ID:           userID,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		Password:     user.Password,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		ModifiedAt:   user.ModifiedAt,
	}
}

func TestCreateUser(t *testing.T) {
	coll := getUserCollection(t)

	t.Run("Creates a new user and throws an error when the same user should be created again", func(t *testing.T) {
		user := createRandomUser(t, coll)

		_, err := coll.CreateUser(context.Background(), user)
		require.Error(t, err) // Duplicate user error
	})
}

func TestGetAllUsers(t *testing.T) {
	coll := getUserCollection(t)

	for i := 0; i < 10; i++ {
		_ = createRandomUser(t, coll)
	}

	pagination := Pagination{
		PageID:   1,
		PageSize: 5,
	}

	t.Run("Gets all users with pagination", func(t *testing.T) {
		ctx := context.Background()
		users, err := coll.GetAllUsers(ctx, pagination)
		require.NoError(t, err)
		require.NotEmpty(t, users)

		require.Equal(t, int(pagination.PageSize), len(users))
	})
}

func TestGetUserByID(t *testing.T) {
	coll := getUserCollection(t)

	createdUser := createRandomUser(t, coll)

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
			gotUser, err := coll.GetUserByID(context.Background(), tc.userID)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedUser.Email, gotUser.Email)
			require.NoError(t, util.CheckPassword(tc.expectedUser.Password, gotUser.PasswordHash))
			require.Equal(t, tc.expectedUser.Role, gotUser.Role)
			require.WithinDuration(t, time.Unix(tc.expectedUser.CreatedAt, 0), time.Unix(gotUser.CreatedAt, 0), 1*time.Second)
			require.WithinDuration(t, time.Unix(tc.expectedUser.ModifiedAt, 0), time.Unix(gotUser.ModifiedAt, 0), 1*time.Second)
		})
	}
}

func TestUpdateUserByID(t *testing.T) {
	coll := getUserCollection(t)

	createdUser := createRandomUser(t, coll)

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
			modifiedCount, err := coll.UpdateUserByID(context.Background(), tc.userID, tc.userUpdate)
			require.Equal(t, tc.modifiedCount, modifiedCount)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if modifiedCount < 1 {
				return
			}

			updatedUser, err := coll.GetUserByID(context.Background(), tc.userID)
			require.NoError(t, err)

			expectedUser := User{
				ID:         createdUser.ID,
				Email:      userUpdate.Email,
				Password:   userUpdate.Password,
				Role:       createdUser.Role,
				CreatedAt:  createdUser.CreatedAt,
				ModifiedAt: time.Now().Unix(),
			}

			require.Equal(t, expectedUser.ID, updatedUser.ID)
			require.Equal(t, expectedUser.Email, updatedUser.Email)
			require.NoError(t, util.CheckPassword(expectedUser.Password, updatedUser.PasswordHash))
			require.Equal(t, expectedUser.Role, updatedUser.Role)
			require.WithinDuration(t, time.Unix(expectedUser.CreatedAt, 0), time.Unix(updatedUser.CreatedAt, 0), 1*time.Second)
			require.WithinDuration(t, time.Unix(expectedUser.ModifiedAt, 0), time.Unix(updatedUser.ModifiedAt, 0), 1*time.Second)
		})
	}
}

func TestDeleteUserByID(t *testing.T) {
	coll := getUserCollection(t)

	createdUser := createRandomUser(t, coll)

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
			deleteCount, err := coll.DeleteUserByID(context.Background(), tc.userID)
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
