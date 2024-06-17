package main

import (
	"gitlab.com/matchmaker/config"
	"gitlab.com/matchmaker/internal/server"
)

func main() {
	cfg := config.New()
	server.Start(cfg)
}
