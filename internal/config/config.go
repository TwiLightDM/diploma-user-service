package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	GRPCPort string

	JWT struct {
		Secret          string
		AccessDuration  time.Duration
		RefreshDuration time.Duration
	}

	SaltLength int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env file not found")
	}

	cfg := &Config{}

	cfg.DB.Host = os.Getenv("POSTGRES_HOST")
	cfg.DB.Port = os.Getenv("POSTGRES_PORT")
	cfg.DB.User = os.Getenv("POSTGRES_USER")
	cfg.DB.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.DB.Name = os.Getenv("POSTGRES_DB")

	cfg.GRPCPort = os.Getenv("GRPC_PORT")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	access, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		log.Fatal("invalid ACCESS_TOKEN_DURATION")
	}
	cfg.JWT.AccessDuration = access

	refresh, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		log.Fatal("invalid REFRESH_TOKEN_DURATION")
	}
	cfg.JWT.RefreshDuration = refresh

	salt, err := strconv.Atoi(os.Getenv("SALT_LENGTH"))
	if err != nil {
		log.Fatal("invalid SALT_LENGTH")
	}
	cfg.SaltLength = salt

	return cfg
}
