package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	Port        string
	DatabaseURL string
	CORSOrigins string
}

func Load() Config {
	loadDotEnv()
	return Config{
		Host:        env("BACKEND_HOST", "0.0.0.0"),
		Port:        env("BACKEND_PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/nisst?sslmode=disable"),
		CORSOrigins: env("CORS_ORIGINS", "http://localhost:3000"),
	}
}

func loadDotEnv() {
	candidates := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}

	wd, err := os.Getwd()
	if err == nil {
		dir := wd
		for i := 0; i < 6; i++ {
			candidates = append(candidates, filepath.Join(dir, ".env"))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	seen := map[string]bool{}
	for _, p := range candidates {
		clean := filepath.Clean(p)
		if seen[clean] {
			continue
		}
		seen[clean] = true
		_ = godotenv.Load(clean)
	}
}

func env(k, v string) string {
	if s := os.Getenv(k); s != "" {
		return s
	}
	return v
}
