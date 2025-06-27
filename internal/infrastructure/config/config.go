package config

import (
	"encoding/json"
	"errors"
	"os"
)

func Read(configPath string) (Config, error) {
	var config Config
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.New("failed to open config file: " + err.Error())
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, errors.New("failed to parse config file: " + err.Error())
	}
	return config, nil
}
