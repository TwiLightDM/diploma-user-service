package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
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
		log.Fatal(".env file didn't found")
	}

	cfg := &Config{}

	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.Name = os.Getenv("POSTGRES_DB")

	cfg.Redis.Host = os.Getenv("REDIS_HOST")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

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
