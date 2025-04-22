package util

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ConnectionString string `mapstructure:"CONNECTION_STRING"`
}

func LoadConfig(path string) *Config {
	viper.SetConfigType("dotenv")
	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	config := &Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return config
}
