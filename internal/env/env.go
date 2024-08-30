package env

import (
	"github.com/Cirqach/GoStream/cmd/logger"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(logger.GetFuncName(0), "Error loading .env file")
	}
}
