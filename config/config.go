package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	config map[string]string
}

func NewConfig() *Config {
	env, err := godotenv.Read()
	if err != nil {
		panic("env not loaded")
	}
	return &Config{
		config: env,
	}
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.config["PG_HOST"],
		c.config["PG_USER"],
		c.config["PG_PASSWORD"],
		c.config["PG_DB_NAME"],
		c.config["PG_PORT"],
	)
}
