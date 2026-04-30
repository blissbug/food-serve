package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)
type Config struct{
  DBUser string 
	DBPass string
	DBName string 
	DBAddr string 
}

func LoadConfig()(Config){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return Config{
		DBUser: getEnv("DBUser","root"),
    DBName: getEnv("DBName","foodserve"),
		DBPass: getEnv("DBPass",""),
		DBAddr: getEnv("DBAddr","127.0.0.1:3306"),
	}
}

func getEnv(key string,fallback string)string{
  value := os.Getenv(key)
	if value == ""{
		return fallback
	}
	return value
}