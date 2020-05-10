package configs

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

const appName = "gif_bot"

type Config struct {
	ApiToken string `required:"true"`
	Debug    bool   `default:"true"`
	Dsn      string `default:"localhost"`
}

func NewConfig() *Config {

	var c Config
	err := envconfig.Process(appName, &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &c
}
