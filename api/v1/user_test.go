package api

import (
	"net/http/httptest"
	"testing"

	mock_db "github.com/PfMartin/wegonice-api/db/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestUnitRegisterUser(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store mock_db.MockDBStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.name, "Success")
	}
}
