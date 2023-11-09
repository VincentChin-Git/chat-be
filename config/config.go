package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	RedisAddr       string
	RedisPass       string
	TokenSecret     string
	AwsAccessKey    string
	AwsSecretKey    string
	AwsBucketName   string
	AwsBucketFolder string
	AwsBucketPrefix string
}

func GetConfig() Config {
	err := godotenv.Load(".env.local")

	if err != nil {
		fmt.Println("No ENV File")
	}

	databaseURL := os.Getenv("DatabaseURL")
	redisAddr := os.Getenv("RedisAddr")
	redisPass := os.Getenv("RedisPass")
	tokenSecret := os.Getenv("TokenSecret")
	awsAccessKey := os.Getenv("AwsAccessKey")
	awsSecretKey := os.Getenv("AwsSecretKey")
	awsBucketName := os.Getenv("AwsBucketName")
	awsBucketFolder := os.Getenv("AwsBucketFolder")
	awsBucketPrefix := os.Getenv("AwsBucketPrefix")

	temp := Config{
		DatabaseURL:     databaseURL,
		RedisAddr:       redisAddr,
		RedisPass:       redisPass,
		TokenSecret:     tokenSecret,
		AwsAccessKey:    awsAccessKey,
		AwsSecretKey:    awsSecretKey,
		AwsBucketName:   awsBucketName,
		AwsBucketFolder: awsBucketFolder,
		AwsBucketPrefix: awsBucketPrefix,
	}

	return temp
}
