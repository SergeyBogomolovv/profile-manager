package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort    int    `mapstructure:"grpc_port"`
	PostgresURL string `mapstructure:"postgres_url"`
}

func MustLoadConfig(path string) *Config {
	viper.SetConfigFile(path)

	viper.BindEnv("postgres_url", "POSTGRES_URL")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("could not unmarshal config file: %v", err)
		return nil
	}

	return &cfg
}
