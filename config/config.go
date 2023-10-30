package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	RedisAddr   string
	RedisPass   string
	TokenSecret string
}

func GetConfig() Config {
	err := godotenv.Load(".env.local")

	if err != nil {
		fmt.Println("No ENV File")
	}

	databaseURL := os.Getenv("DatabaseURL")
	databaseName := os.Getenv("DatabaseName")
	redisAddr := os.Getenv("RedisAddr")
	redisPass := os.Getenv("RedisPass")
	tokenSecret := os.Getenv("TokenSecret")

	fmt.Println(databaseName, databaseURL)
	temp := Config{DatabaseURL: databaseURL, RedisAddr: redisAddr, RedisPass: redisPass, TokenSecret: tokenSecret}

	return temp
}
