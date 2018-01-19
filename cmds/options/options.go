package options

import (
	"log"
	"os"

	"github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/util"
)

type Config struct {
	Provider string
}

const (
	keyProvider = "provider"
	envProvider = "PROVIDER"
)

func NewConfig() *Config {
	var err error
	provider := os.Getenv(envProvider)
	if provider == "" {
		provider, err = util.ReadSecretKeyFromFile(cloud.SecretDefaultLocation, keyProvider)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return &Config{
		Provider: provider,
	}
}
