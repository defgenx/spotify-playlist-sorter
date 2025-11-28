package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Server  ServerConfig
	Spotify SpotifyConfig
	Session SessionConfig
}

type ServerConfig struct {
	Port         string `env:"PORT" envDefault:"3001"`
	FrontendURL  string `env:"FRONTEND_URL" envDefault:"http://localhost:3000"`
	AllowOrigins string `env:"CORS_ORIGINS" envDefault:"http://localhost:3000"`
}

type SpotifyConfig struct {
	ClientID     string `env:"SPOTIFY_CLIENT_ID,required"`
	ClientSecret string `env:"SPOTIFY_CLIENT_SECRET,required"`
	RedirectURL  string `env:"SPOTIFY_REDIRECT_URL" envDefault:"http://localhost:3001/api/auth/callback"`
}

type SessionConfig struct {
	Secret string `env:"SESSION_SECRET,required"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
