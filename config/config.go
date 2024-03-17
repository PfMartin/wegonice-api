package config

import (
	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

type Config struct {
	DBName     string `mapstructure:"WEGONICE_DB"`
	DBUser     string `mapstructure:"WEGONICE_USER"`
	DBPassword string `mapstructure:"WEGONICE_PWD"`
	DBURI      string `mapstructure:"WEGONICE_URI"`
	APIURL     string `mapstructure:"API_URL"`
}

func NewConfig(configPath string, configName string) (config Config, err error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		log.Fatal().Msgf("failed to read config: %s", err)
		return
	}

	err = viper.Unmarshal(&config)

	return config, err
}
