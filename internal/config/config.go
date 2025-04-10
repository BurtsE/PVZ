package config

import (
	"log"
	"os"
)

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetApplicationPort() string {
	port := getEnvironmentValue("APPLICATION_PORT")
	return port
}
func GetSecretKey() string {
	return getEnvironmentValue("SECRET_KEY")
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}
	return os.Getenv(key)
}
