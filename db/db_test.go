package db

import (
	"os"
	"testing"

	"github.com/PfMartin/wegonice-api/config"
	"github.com/stretchr/testify/require"
)

func TestConnectToDatabase(t *testing.T) {
	t.Run("Connects to database", func(t *testing.T) {
		conf := config.Config{
			DBName:     os.Getenv("WEGONICE_DB"),
			DBUser:     os.Getenv("WEGONICE_USER"),
			DBPassword: os.Getenv("WEGONICE_PWD"),
			DBURI:      os.Getenv("WEGONICE_URI"),
		}

		if conf.DBName == "" || conf.DBUser == "" || conf.DBPassword == "" || conf.DBURI == "" {
			viperConf, err := config.NewConfig("./", "test.env")
			require.NoError(t, err)

			conf = viperConf
		}

		client, cancel := NewDatabase(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
		defer cancel()
		require.NotNil(t, client)
	})

}
