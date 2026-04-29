package config

import (
	"fmt"
	"os"
)

// GetEnv returns the value of an environment variable or a fallback default.
func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// MustEnv returns the value of a required environment variable.
// It panics if the variable is not set, making misconfiguration immediately visible at startup.
func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return v
}
