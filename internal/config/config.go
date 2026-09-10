package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseDSN    string
	Port           string
	AllowedOrigins []string
}

func Load() *Config {
	return &Config{
		DatabaseDSN:    getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=caesar_quiz port=5432 sslmode=disable"),
		Port:           getEnv("PORT", "8080"),
		AllowedOrigins: parseOrigins(getEnv("CORS_ORIGINS", "*")),
	}
}

// parseOrigins memecah daftar origin dipisah koma. Default "*" (semua origin)
// supaya mudah di-deploy; set CORS_ORIGINS untuk membatasi ke domain produksi.
func parseOrigins(raw string) []string {
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			origins = append(origins, p)
		}
	}
	if len(origins) == 0 {
		origins = []string{"*"}
	}
	return origins
}

func getEnv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
