package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	App struct {
		Port string `env:"PORT"`
	}
	OAuth struct {
		ClientId     string `env:"OAUTH_CLIENT_ID"`
		ClientSecret string `env:"OAUTH_CLIENT_SECRET"`
		RedirectUrl  string `env:"OAUTH_REDIRECT_URL"`
	}
	Pusher struct {
		InstanceId string `env:"PUSHER_INSTANCE_ID"`
		PrivateKey string `env:"PUSHER_PRIVATE_KEY"`
	}
}

var (
	c    Config
	once sync.Once
)

func Parser() Config {
	once.Do(func() {
		if os.Getenv("PORT") == "" {
			if err := godotenv.Load(); err != nil {
				panic(err)
			}
		}
		if err := env.Parse(&c); err != nil {
			panic(err)
		}
	})
	return c
}
