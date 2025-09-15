package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	DSN          string
	JWTSecret    string
	JWTAccessTTL time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	ssl := os.Getenv("DB_SSLMODE")

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	jwtSecret := mustEnv("JWT_SECRET")
	jwtAccessTTL := mustDuration("JWT_ACCESS_TTL", "15m")

	dsn := "host=" + host +
		" user=" + user +
		" password=" + pass +
		" dbname=" + name +
		" port=" + port +
		" sslmode=" + ssl

	cfg := &Config{
		AppPort:      appPort,
		DSN:          dsn,
		JWTSecret:    jwtSecret,
		JWTAccessTTL: jwtAccessTTL,
	}

	log.Printf("config loaded")
	return cfg
}

// helper ambil env wajib
func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env: %s", key)
	}
	return val
}

// helper parse durasi (contoh: "15m", "1h")
func mustDuration(key, def string) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		val = def
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("invalid duration for %s: %v", key, err)
	}
	return d
}
