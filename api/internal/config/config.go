package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	Port              string
	AdminUser         string
	AdminPassword     string
	WhatsAppProvider  string
	EvolutionAPIURL   string
	EvolutionAPIKey   string
	EvolutionInstance string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/barber?sslmode=disable"),
		Port:              getEnv("PORT", "8080"),
		AdminUser:         getEnv("ADMIN_USER", "admin"),
		AdminPassword:     getEnv("ADMIN_PASSWORD", "admin123"),
		WhatsAppProvider:  getEnv("WHATSAPP_PROVIDER", "evolution"),
		EvolutionAPIURL:   getEnv("EVOLUTION_API_URL", "http://localhost:8085"),
		EvolutionAPIKey:   getEnv("EVOLUTION_API_KEY", ""),
		EvolutionInstance: getEnv("EVOLUTION_INSTANCE", "barber"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
