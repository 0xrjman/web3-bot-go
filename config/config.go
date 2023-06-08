package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	EthereumURL   string
	KEYSTORE_PATH string
	SIGNER        string
	EthereumPK    string
}

func LoadConfig() *Config {
	cfg := &Config{
		DBHost:        getEnv("DB_HOST", ""),
		DBPort:        getEnv("DB_PORT", ""),
		DBUser:        getEnv("DB_USER", ""),
		DBPassword:    getEnv("DB_PWD", ""),
		DBName:        getEnv("DB_NAME", ""),
		EthereumURL:   getEnv("ETH_URL", ""),
		EthereumPK:    getEnv("ETH_PK", ""),
		KEYSTORE_PATH: getEnv("KEYSTORE_PATH", ""),
		SIGNER:        getEnv("SIGNER", ""),
	}

	return cfg
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue != "" {
			return defaultValue
		}
		log.Fatalf("Environment variable %s is required", key)
	}
	return value
}
