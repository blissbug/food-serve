package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser string
	DBPass string
	DBName string
	DBAddr string

	RedisAddr string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return Config{
		DBUser:    getEnv("DBUser", "root"),
		DBName:    getEnv("DBName", "foodserve"),
		DBPass:    getEnv("DBPass", ""),
		DBAddr:    getEnv("DBAddr", "127.0.0.1:3306"),
		RedisAddr: getEnv("RedisAddr", "127.0.0.1:6379"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
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
	return nil
}
