package config

import (
	"os"
)

type Config struct {
	DatabaseURL      string
	JWTSecret       string
	OpenAIAPIKey    string
	StripeSecretKey string
	StripeWebhookSecret string
	CORSOrigins     string
	Port           string
}

func Load() *Config {
	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://localhost/symbol_quest?sslmode=disable"),
		JWTSecret:       getEnv("JWT_SECRET", "your-256-bit-secret"),
		OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),
		StripeSecretKey: getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		CORSOrigins:     getEnv("CORS_ORIGINS", "http://localhost:5173,https://symbol-quest.vercel.app"),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}