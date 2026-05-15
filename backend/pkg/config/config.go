package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBUser string `mapstructure:"DBUser"`
	DBPass string `mapstructure:"DBPass"`
	DBName string `mapstructure:"DBName"`
	DBAddr string `mapstructure:"DBAddr"`

	RedisAddr string `mapstructure:"RedisAddr"`
	JWTSecret string `mapstructure:"JWTSecret"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		return cfg, fmt.Errorf("No .env file found")
	}

	err = viper.Unmarshal(&cfg)

	if err != nil {
		return cfg, fmt.Errorf("Error reading .env file: %s", err)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.DBName == "" {
		return fmt.Errorf("DBName is required")
	}
	if c.DBPass == "" {
		return fmt.Errorf("DBPass is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DBUser is required")
	}
	if c.DBAddr == "" {
		return fmt.Errorf("DBAddr is required")
	}
	if c.RedisAddr == "" {
		return fmt.Errorf("RedisAddr is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWTSecret is required")
	}
	return nil
}
