package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitGetConfig(t *testing.T) {
	t.Run("Gets correct config", func(t *testing.T) {
		conf, err := NewConfig("./", "test.env")
		require.NoError(t, err)

		require.Equal(t, "wegonice", conf.DBName)
		require.Equal(t, "niceUser", conf.DBUser)
		require.Equal(t, "nicePassword", conf.DBPassword)
		require.Equal(t, "mongodb://localhost:27017", conf.DBURI)
	})
}
