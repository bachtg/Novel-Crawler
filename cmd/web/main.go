package main

import (
	"go.uber.org/zap"
	"novel_crawler/router"

	"novel_crawler/config"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		config.Cfg.Logger.Fatal("could not load config", zap.Error(err))
	}
	router.Start()
}
