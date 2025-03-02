package config

import (
	"log"

	"github.com/spf13/viper"
)

type (
	Config struct {
		GRPC        GRPC   `mapstructure:"grpc"`
		HTTP        HTTP   `mapstructure:"http"`
		PostgresURL string `mapstructure:"postgres_url"`
		RedisURL    string `mapstructure:"redis_url"`
		OAuth       OAuth  `mapstructure:"oauth"`
		JWT         JWT    `mapstructure:"jwt"`
	}
	GRPC struct {
		Port int `mapstructure:"port"`
	}
	HTTP struct {
		Port int `mapstructure:"port"`
	}
	OAuth struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
	}
	JWT struct {
		SecretKey string `mapstructure:"secret_key"`
	}
)

func MustLoadConfig(path string) *Config {
	viper.SetConfigFile(path)

	viper.BindEnv("jwt.secret_key", "JWT_SECRET_KEY")
	viper.BindEnv("redis_url", "REDIS_URL")
	viper.BindEnv("postgres_url", "POSTGRES_URL")
	viper.BindEnv("oauth.client_id", "OAUTH_CLIENT_ID")
	viper.BindEnv("oauth.client_secret", "OAUTH_CLIENT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("could not read config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("could not unmarshal config file: %v", err)
		return nil
	}

	return &cfg
}
