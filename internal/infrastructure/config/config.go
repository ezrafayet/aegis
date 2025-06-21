package config

import (
	"encoding/json"
	"errors"
	"os"
	"othnx/internal/domain"
)

func Read(configPath string) (domain.Config, error) {
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

	// err = domain.ValidateConfig(config)
	// if err != nil {
	// 	return domain.Config{}, errors.New("invalid config file: " + err.Error())
	// }

	return config, nil
}
