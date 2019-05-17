package util

import (
	"os"
)

var DEFAULT_API_BASE_URL = "https://api.whiteblock.io"

var ApiEnvVar = "API_URL"

// Get API_BASE_URL from the environment variable API_URL or fallback to the default
func GetEnvVar(envVar, fallback string) string {
	val, exists := os.LookupEnv(envVar)
	if exists {
		return val
	}
	return fallback
}

var ApiBaseURL = GetEnvVar(ApiEnvVar, DEFAULT_API_BASE_URL)
