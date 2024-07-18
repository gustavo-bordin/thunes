package config

import (
	"log"

	"github.com/spf13/viper"
)

type ApiConfig struct {
	Mongo MongoConfig
	Port  int16
}

func NewApiConfig() ApiConfig {
	return ApiConfig{}
}

func (c *ApiConfig) Load() {
	viper.AddConfigPath(".")
	viper.SetConfigName("api")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}
}
