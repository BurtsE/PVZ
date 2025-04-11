package config

import (
	"fmt"
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

func GetUsersPostgresURL() string {
	user := getEnvironmentValue("POSTGRES_USER")
	password := getEnvironmentValue("POSTGRES_PASSWORD")
	database := getEnvironmentValue("POSTGRES_DB")
	//port := getEnvironmentValue("POSTGRES_PORT")
	host := getEnvironmentValue("DATABASE_HOST")
	return fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", host, user, database, password)
}
