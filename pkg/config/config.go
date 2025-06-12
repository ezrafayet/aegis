package config

import (
	"encoding/json"
	"errors"
	"os"
	"othnx/internal/domain"
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
		return domain.Config{}, errors.New("failed to parse config file: " + err.Error())
	}

	err = validateConfig(config)
	if err != nil {
		return domain.Config{}, errors.New("invalid config file: " + err.Error())
	}

	return config, nil
}

func validateConfig(config domain.Config) error {
	// todo
	return nil
}
