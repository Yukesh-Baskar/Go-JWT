package utils

import (
	config "go-jwt/configs"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

func GetConfig() (port, mongoDBUri string, err error) {
	data, err := os.ReadFile("./configs/config.yaml")
	if err != nil {
		return "", "", err
	}
	cfg := config.Config{}
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return "", "", err
	}

	return cfg.Port, cfg.MonogDBUri, nil
}

func GetHashedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}