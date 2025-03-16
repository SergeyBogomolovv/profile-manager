package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort    int    `mapstructure:"grpc_port"`
	PostgresURL string `mapstructure:"postgres_url"`
	RabbitmqURL string `mapstructure:"rabbitmq_url"`
	JwtSecret   string `mapstructure:"jwt_secret"`
}

func MustLoadConfig(path string) *Config {
	viper.SetConfigFile(path)

	viper.BindEnv("postgres_url", "POSTGRES_URL")
	viper.BindEnv("rabbitmq_url", "RABBITMQ_URL")
	viper.BindEnv("jwt_secret", "JWT_SECRET")

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
