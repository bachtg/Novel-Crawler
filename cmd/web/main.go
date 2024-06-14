package main

import (
	"go.uber.org/zap"

	"novel_crawler/config"
	"novel_crawler/router"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		config.Cfg.Logger.Fatal("could not load config", zap.Error(err))
	}
	router.Start()
}
