package db

import (
	"testing"

	"github.com/PfMartin/wegonice-api/config"
	"github.com/stretchr/testify/require"
)

func TestConnectToDatabase(t *testing.T) {
	t.Run("Connects to database", func(t *testing.T) {
		conf, err := config.NewConfig("./", "test.env")
		require.NoError(t, err)

		client, cancel := NewDatabase(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
		defer cancel()
		require.NotNil(t, client)
	})

}
