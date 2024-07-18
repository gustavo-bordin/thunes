package config

import (
	"log"

	"github.com/spf13/viper"
)

type NgrokConfig struct {
	Url string
}

type ThunesConfig struct {
	Username string
	Password string
	HostUrl  string
}

type CliConfig struct {
	Mongo  MongoConfig
	Thunes ThunesConfig
	Ngrok  NgrokConfig
}

func NewCliConfig() CliConfig {
	return CliConfig{Mongo: MongoConfig{}}
}

func (c *CliConfig) Load() {
	viper.AddConfigPath(".")
	viper.SetConfigName("cli")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}
}
