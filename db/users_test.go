package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T, store *MongoDBStore) User {
	t.Helper()

	user := User{
		Email:      util.RandomEmail(),
		Password:   util.RandomString(6),
		Role:       UserRole,
		IsActive:   false,
		CreatedAt:  time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	}

	ctx := context.Background()
	insertedID, err := store.CreateUser(ctx, user)

	require.NoError(t, err)
	require.False(t, insertedID.IsZero())

	userID := insertedID.Hex()

	hashedPassword, err := util.HashPassword(user.Password)
	require.NoError(t, err)

	fmt.Println(user.CreatedAt)

	return User{
		ID:           userID,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		Password:     user.Password,
		IsActive:     user.IsActive,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		ModifiedAt:   user.ModifiedAt,
	}
}

func TestUnitCreateUser(t *testing.T) {
	store := getMongoDBStore(t)

	t.Run("Creates a new user and throws an error when the same user should be created again", func(t *testing.T) {
		user := createRandomUser(t, store)

		_, err := store.CreateUser(context.Background(), user)
		require.Error(t, err) // Duplicate user error
	})
}

func TestUnitGetAllUsers(t *testing.T) {
	store := getMongoDBStore(t)

	for i := 0; i < 10; i++ {
		_ = createRandomUser(t, store)
	}

	pagination := Pagination{
		PageID:   1,
		PageSize: 5,
	}

	t.Run("Gets all users with pagination", func(t *testing.T) {
		ctx := context.Background()
		users, err := store.GetAllUsers(ctx, pagination)
		require.NoError(t, err)
		require.NotEmpty(t, users)

		require.Equal(t, int(pagination.PageSize), len(users))
	})
}

func TestUnitGetUserByID(t *testing.T) {
	store := getMongoDBStore(t)

	createdUser := createRandomUser(t, store)

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
			gotUser, err := store.GetUserByID(context.Background(), tc.userID)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedUser.Email, gotUser.Email)
			require.NoError(t, util.CheckPassword(tc.expectedUser.Password, gotUser.PasswordHash))
			require.Equal(t, tc.expectedUser.Role, gotUser.Role)
			require.Equal(t, tc.expectedUser.IsActive, gotUser.IsActive)
			require.WithinDuration(t, time.Unix(tc.expectedUser.CreatedAt, 0), time.Unix(gotUser.CreatedAt, 0), 5*time.Second)
			require.WithinDuration(t, time.Unix(tc.expectedUser.ModifiedAt, 0), time.Unix(gotUser.ModifiedAt, 0), 5*time.Second)
		})
	}
}

func TestUnitGetUserByEmail(t *testing.T) {
	store := getMongoDBStore(t)

	createdUser := createRandomUser(t, store)

	testCases := []struct {
		name         string
		email        string
		hasError     bool
		expectedUser User
	}{
		{
			name:         "Success",
			email:        createdUser.Email,
			hasError:     false,
			expectedUser: createdUser,
		},
		{
			name:     "Fail with email not found",
			email:    "test@email.com",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotUser, err := store.GetUserByEmail(context.Background(), tc.email)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedUser.Email, gotUser.Email)
			require.NoError(t, util.CheckPassword(tc.expectedUser.Password, gotUser.PasswordHash))
			require.Equal(t, tc.expectedUser.Role, gotUser.Role)
			require.Equal(t, tc.expectedUser.IsActive, gotUser.IsActive)
			require.WithinDuration(t, time.Unix(tc.expectedUser.CreatedAt, 0), time.Unix(gotUser.CreatedAt, 0), 5*time.Second)
			require.WithinDuration(t, time.Unix(tc.expectedUser.ModifiedAt, 0), time.Unix(gotUser.ModifiedAt, 0), 5*time.Second)
		})
	}
}

func TestUnitUpdateUserByID(t *testing.T) {
	store := getMongoDBStore(t)

	createdUser := createRandomUser(t, store)

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
			modifiedCount, err := store.UpdateUserByID(context.Background(), tc.userID, tc.userUpdate)
			require.Equal(t, tc.modifiedCount, modifiedCount)

			if tc.hasError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if modifiedCount < 1 {
				return
			}

			updatedUser, err := store.GetUserByID(context.Background(), tc.userID)
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
			require.Equal(t, expectedUser.IsActive, updatedUser.IsActive)
			require.WithinDuration(t, time.Unix(expectedUser.CreatedAt, 0), time.Unix(updatedUser.CreatedAt, 0), 5*time.Second)
			require.WithinDuration(t, time.Unix(expectedUser.ModifiedAt, 0), time.Unix(updatedUser.ModifiedAt, 0), 5*time.Second)
		})
	}
}

func TestUnitDeleteUserByID(t *testing.T) {
	store := getMongoDBStore(t)

	createdUser := createRandomUser(t, store)
	authorUser := createRandomUser(t, store)
	recipeUser := createRandomUser(t, store)

	recipeAuthor := createRandomAuthor(t, store, recipeUser.ID)

	testCases := []struct {
		name                 string
		userID               string
		hasError             bool
		hasReferencingRecipe bool
		hasReferencingAuthor bool
		deleteCount          int64
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
		{
			name:                 "Fail due to referencing author",
			userID:               authorUser.ID,
			hasError:             true,
			hasReferencingAuthor: true,
			deleteCount:          0,
		},
		{
			name:                 "Fail due to referencing recipe",
			userID:               recipeUser.ID,
			hasError:             true,
			hasReferencingRecipe: true,
			deleteCount:          0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.hasReferencingAuthor {
				createRandomAuthor(t, store, authorUser.ID)
			}

			if tc.hasReferencingRecipe {
				createRandomRecipe(t, store, recipeUser.ID, recipeAuthor.ID)
			}

			deleteCount, err := store.DeleteUserByID(context.Background(), tc.userID)
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
