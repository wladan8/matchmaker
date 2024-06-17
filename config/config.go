package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	ServerPort      int     `env:"SERVER_PORT"`
	GroupSize       int     `env:"GROUP_SIZE"`
	DiffSkill       float64 `env:"DIFF_SKILL" `
	DiffLatency     float64 `env:"DIFF_LATENCY"`
	TickerFrequency int     `env:"TICKER_FREQUENCY"`
}

func New() *ServerConfig {
	cfg := &ServerConfig{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
