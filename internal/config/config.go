package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	// PgDBUri    string `mapstructure:"POSTGRES_LOCAL_URI"`
	Port string `mapstructure:"PORT"`

	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
