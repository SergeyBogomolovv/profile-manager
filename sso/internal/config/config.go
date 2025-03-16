package config

import (
	"log"

	"github.com/spf13/viper"
)

type (
	Config struct {
		GrpcPort    int    `mapstructure:"grpc_port"`
		HttpPort    int    `mapstructure:"http_port"`
		PostgresURL string `mapstructure:"postgres_url"`
		RedisURL    string `mapstructure:"redis_url"`
		RabbitmqURL string `mapstructure:"rabbitmq_url"`
		JwtSecret   string `mapstructure:"jwt_secret"`
		OAuth       OAuth  `mapstructure:"oauth"`
	}
	OAuth struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
	}
)

func MustLoadConfig(path string) *Config {
	viper.SetConfigFile(path)

	viper.BindEnv("jwt_secret", "JWT_SECRET")
	viper.BindEnv("redis_url", "REDIS_URL")
	viper.BindEnv("postgres_url", "POSTGRES_URL")
	viper.BindEnv("oauth.client_id", "OAUTH_CLIENT_ID")
	viper.BindEnv("oauth.client_secret", "OAUTH_CLIENT_SECRET")
	viper.BindEnv("rabbitmq_url", "RABBITMQ_URL")

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
