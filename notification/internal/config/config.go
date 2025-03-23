package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort      int    `mapstructure:"grpc_port"`
	PostgresURL   string `mapstructure:"postgres_url"`
	RedisURL      string `mapstructure:"redis_url"`
	RabbitmqURL   string `mapstructure:"rabbitmq_url"`
	TelegramToken string `mapstructure:"telegram_token"`
	SMTP          SMTP   `mapstructure:"smtp"`
	JwtSecret     string `mapstructure:"jwt_secret"`
}

type SMTP struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"password"`
}

func MustLoadConfig(path string) *Config {
	viper.SetConfigFile(path)

	viper.BindEnv("postgres_url", "POSTGRES_URL")
	viper.BindEnv("redis_url", "REDIS_URL")
	viper.BindEnv("rabbitmq_url", "RABBITMQ_URL")
	viper.BindEnv("telegram_token", "TELEGRAM_TOKEN")
	viper.BindEnv("smtp.password", "SMTP_PASSWORD")
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
