package config

import (
	"log"

	"github.com/spf13/viper"
)

type (
	Config struct {
		GRPC     GRPC     `mapstructure:"grpc"`
		HTTP     HTTP     `mapstructure:"http"`
		Database Database `mapstructure:"database"`
		OAuth    OAuth    `mapstructure:"oauth"`
	}
	GRPC struct {
		Port int `mapstructure:"port"`
	}
	HTTP struct {
		Port int `mapstructure:"port"`
	}
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
	}
	OAuth struct {
		ClientID     string `mapstructure:"client_id"`
		ClientSecret string `mapstructure:"client_secret"`
		RedirectURL  string `mapstructure:"redirect_url"`
	}
)

func MustLoadConfig(path string) *Config {
	viper.SetConfigFile(path)

	viper.BindEnv("database.password", "DB_PASSWORD")
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
