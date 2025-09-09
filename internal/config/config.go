package config

import "os"

type Config struct {
	DatabaseDSN string
	Port        string
	MaxPlayers  int
}

func Load() *Config {
	return &Config{
		DatabaseDSN: getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=caesar_quiz port=5432 sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
