package config

import (
	"os"
)

var (
	// Config store current configuration values.
	Config = loadConfig()
)

type config struct {
	DB_HOST        string
	DB_NAME        string
	DB_USER        string
	DB_PASS        string
	APP_JWT_SECRET string
}

func loadConfig() *config {
	c := new(config)

	c.DB_HOST = os.Getenv("DB_HOST")
	c.DB_NAME = os.Getenv("DB_NAME")
	c.DB_USER = os.Getenv("DB_USER")
	c.DB_PASS = os.Getenv("DB_PASS")
	c.APP_JWT_SECRET = os.Getenv("APP_JWT_SECRET")

	return c
}
