package config

import (
	"log"
	"vendors/pkg/lib/utils"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string  `yaml:"env"`
	Server  Server  `yaml:"server"`
	MongoDB MongoDB `yaml:"mongodb"`
}

type Server struct {
	Address string `yaml:"address"`
}

type MongoDB struct {
	URI               string `yaml:"uri"`
	Database          string `yaml:"database"`
	CinemaCollection  string `yaml:"cinemaCollection"`
	TheatreCollection string `yaml:"theatreCollection"`
	FoodCollection    string `yaml:"foodCollection"`
}

func LoadConfig() *Config {
	configPath := "./config/config.yaml"

	if configPath == "" {
		log.Fatalf("config path is not set or config file does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %v", utils.Err(err))
	}

	return &cfg
}
