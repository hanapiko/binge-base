package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	GinMode        string
	TMDBAPIKey     string
	OMDBAPIKey     string
	DBPath         string
	AllowedOrigins []string
	TMDBRateLimit  int
	OMDBRateLimit  int
	CacheDuration  int
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		TMDBAPIKey:     getEnv("TMDB_API_KEY", ""),
		OMDBAPIKey:     getEnv("OMDB_API_KEY", ""),
		DBPath:         getEnv("DB_PATH", "./database/bingebase.db"),
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		TMDBRateLimit:  getEnvAsInt("TMDB_RATE_LIMIT", 40),
		OMDBRateLimit:  getEnvAsInt("OMDB_RATE_LIMIT", 1000),
		CacheDuration:  getEnvAsInt("CACHE_DURATION", 3600),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
