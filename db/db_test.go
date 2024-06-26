package db

import (
	"os"
	"testing"

	"github.com/PfMartin/wegonice-api/config"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func getDatabaseConfiguration(t *testing.T) config.Config {
	conf := config.Config{
		DBName:     os.Getenv("WEGONICE_DB"),
		DBUser:     os.Getenv("WEGONICE_USER"),
		DBPassword: os.Getenv("WEGONICE_PWD"),
		DBURI:      os.Getenv("WEGONICE_URI"),
	}

	if conf.DBName == "" || conf.DBUser == "" || conf.DBPassword == "" || conf.DBURI == "" {
		viperConf, err := config.NewConfig("../", ".env")
		require.NoError(t, err)

		conf = viperConf
	}

	return conf
}

func TestUnitConnectToDatabase(t *testing.T) {
	t.Run("Connects to database", func(t *testing.T) {
		conf := getDatabaseConfiguration(t)

		client, cancel := NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
		defer cancel()
		require.NotNil(t, client)
	})
}

func TestUnitGetSortStage(t *testing.T) {
	testCases := []struct {
		name string
		key  string
	}{
		{
			name: "Success with name sorting",
			key:  "name",
		},
		{
			name: "Success with email sorting",
			key:  "email",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedStage := bson.M{"$sort": bson.M{tc.key: 1}}

			require.Equal(t, expectedStage, getSortStage(tc.key))
		},
		)
	}
}
