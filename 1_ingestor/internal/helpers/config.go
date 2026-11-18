package helpers

import (
	"os"
)

// GetEnv
// Returns the environment variable, or defaultValue when it is not set.
func GetEnv(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}
