package config

import (
	"aegix/internal/domain"
	"encoding/json"
	"errors"
	"os"
)

func ReadConfig(configPath string) (domain.Config, error) {
	var config domain.Config

	file, err := os.Open(configPath)
	if err != nil {
		return domain.Config{}, errors.New("failed to open config file: " + err.Error())
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return domain.Config{}, errors.New("failed to parse config: " + err.Error())
	}

	return config, nil
}
