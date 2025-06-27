package config

import (
	"encoding/json"
	"errors"
	"os"
	"othnx/internal/domain/entities"
)

func Read(configPath string) (entities.Config, error) {
	var config entities.Config
	file, err := os.Open(configPath)
	if err != nil {
		return entities.Config{}, errors.New("failed to open config file: " + err.Error())
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return entities.Config{}, errors.New("failed to parse config file: " + err.Error())
	}
	return config, nil
}
