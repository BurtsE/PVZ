package config

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentFunctions(t *testing.T) {
	// Setup and teardown for environment variables
	setEnv := func(key, value string) {
		err := os.Setenv(key, value)
		require.NoError(t, err)
	}
	unsetEnv := func(key string) {
		err := os.Unsetenv(key)
		require.NoError(t, err)
	}

	// Backup and restore the original log.Fatal function
	originalLogFatal := logFatal
	defer func() { logFatal = originalLogFatal }()

	// Mock log.Fatal to prevent tests from exiting
	
	t.Run("GetEnv", func(t *testing.T) {
		t.Run("returns environment value when set", func(t *testing.T) {
			setEnv("ENV", "test")
			defer unsetEnv("ENV")

			result := GetEnv()
			assert.Equal(t, "test", result)
		})

	})

	t.Run("GetApplicationPort", func(t *testing.T) {
		t.Run("returns port when set", func(t *testing.T) {
			setEnv("APPLICATION_PORT", "8080")
			defer unsetEnv("APPLICATION_PORT")

			result := GetApplicationPort()
			assert.Equal(t, "8080", result)
		})

	})

	t.Run("GetSecretKey", func(t *testing.T) {
		t.Run("returns secret key when set", func(t *testing.T) {
			setEnv("SECRET_KEY", "my-secret-key")
			defer unsetEnv("SECRET_KEY")

			result := GetSecretKey()
			assert.Equal(t, "my-secret-key", result)
		})

	})

	t.Run("GetUsersPostgresURL", func(t *testing.T) {
		t.Run("returns correct connection string", func(t *testing.T) {
			setEnv("POSTGRES_USER", "user")
			setEnv("POSTGRES_PASSWORD", "pass")
			setEnv("POSTGRES_DB", "db")
			setEnv("DATABASE_HOST", "localhost")
			defer func() {
				unsetEnv("POSTGRES_USER")
				unsetEnv("POSTGRES_PASSWORD")
				unsetEnv("POSTGRES_DB")
				unsetEnv("DATABASE_HOST")
			}()

			expected := "host=localhost user=user dbname=db password=pass sslmode=disable"
			result := GetUsersPostgresURL()
			assert.Equal(t, expected, result)
		})

	})
}

// logFatal is a variable that holds the log.Fatalf function
// so we can mock it in tests
var logFatal = func(msg string) {
	log.Fatalf(msg)
}
