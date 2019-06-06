package util

import (
	"os"
)

const DefaultAPIBaseURL = "https://api.whiteblock.io"

const ApiEnvVar = "API_URL"

// Get API_BASE_URL from the environment variable API_URL or fallback to the default
func GetEnvVar(envVar, fallback string) string {
	val, exists := os.LookupEnv(envVar)
	if exists {
		return val
	}
	return fallback
}

var ApiBaseURL = GetEnvVar(ApiEnvVar, DefaultAPIBaseURL)
