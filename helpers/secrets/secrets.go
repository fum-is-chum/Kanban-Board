package helpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadSecrets(envPath string, key string) (string, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return "", err
	}

	if val, found := os.LookupEnv("SECRET_JWT"); found {
		return val, nil
	}

	return "", errors.New(fmt.Sprintf("Key %s not found!", key))
}
