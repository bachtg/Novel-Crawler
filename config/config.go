package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger            *zap.Logger
	Address           string `yaml:"address"`
	TruyenFullBaseUrl string `yaml:"truyen_full_base_url"`
	NetTruyenBaseUrl string `yaml:"net_truyen_base_url"`
}

var Cfg Config

func LoadConfig() error {
	// read config from file
	yamlData, err := os.ReadFile("./config.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlData, &Cfg)
	if err != nil {
		return err
	}
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	Cfg.Logger = logger

	Cfg.Logger.Info("Loaded configuration successfully!")

	return err
}
