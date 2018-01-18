package options

import "os"

type Config struct {
	Provider string
}

func NewConfig() *Config {
	return &Config{
		Provider: os.Getenv("PROVIDER"),
	}
}
