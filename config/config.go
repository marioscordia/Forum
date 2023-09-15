package config

import (
	"encoding/json"
	"os"
)

type DB struct {
	Driver string `json:"driver"`
	Dsn string `json:"dsn"`
}

type Config struct {
	Port string `json:"port"`
	DB DB `json:"db"`
}

func NewConfig() (*Config, error) {
	file, err := os.Open("./config/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config *Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}