package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	HttpPort         int    `mapstructure:"http_port"`
	SsoAddr          string `mapstructure:"sso_addr"`
	ProfileAddr      string `mapstructure:"profile_addr"`
	NotificationAddr string `mapstructure:"notification_addr"`
}

func MustLoadConfig(path string) *Config {
	viper.SetConfigFile(path)

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
