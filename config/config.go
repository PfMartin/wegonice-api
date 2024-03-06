package config

import "github.com/spf13/viper"

type Config struct {
	DBName     string `mapstructure:"WEGONICE_DB"`
	DBUser     string `mapstructure:"WEGONICE_USER"`
	DBPassword string `mapstructure:"WEGONICE_PWD"`
	DBURI      string `mapstructure:"WEGONICE_URI"`
}

func NewConfig() (config Config, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
