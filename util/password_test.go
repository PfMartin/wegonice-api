package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUnitPassword(t *testing.T) {
	testCases := []struct {
		name      string
		password1 string
		password2 string
		isEqual   bool
	}{
		{
			name:      "Success with the same password",
			password1: "testpassword",
			password2: "testpassword",
			isEqual:   true,
		},
		{
			name:      "Fail with different passwords",
			password1: "test",
			password2: "password",
			isEqual:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedPassword1, err := HashPassword(tc.password1)
			require.NoError(t, err)
			require.NotEmpty(t, hashedPassword1)

			err = CheckPassword(tc.password2, hashedPassword1)
			if tc.isEqual {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
			}
		})

	}
}
