package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

type ServerConfig struct {
	Address string
	Mode    string
}

type DatabaseConfig struct {
	URL string
}

type LogConfig struct {
	Level  string
	Format string
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/finapp")

	// Set defaults
	viper.SetDefault("server.address", ":8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.url", "postgres://postgres:postgres@localhost:5432/finapp?sslmode=disable")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "console")

	// Environment variable overrides
	viper.AutomaticEnv()
	viper.SetEnvPrefix("FINAPP")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic("Failed to read config: " + err.Error())
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic("Failed to unmarshal config: " + err.Error())
	}

	return &cfg
}