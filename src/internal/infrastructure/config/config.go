package config

import (
	"aegis/internal/domain/entities"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"regexp"
	"strings"
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

	// Replace environment variables in the config
	replaceEnvVars(&config)

	return config, nil
}

func replaceEnvVars(config any) {
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.String:
			if field.CanSet() {
				str := field.String()
				replaced := interpolateEnvVars(str)
				field.SetString(replaced)
			}
		case reflect.Struct:
			if field.CanAddr() {
				replaceEnvVars(field.Addr().Interface())
			}
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.String && elem.CanSet() {
					str := elem.String()
					replaced := interpolateEnvVars(str)
					elem.SetString(replaced)
				}
			}
		}
	}
}

func interpolateEnvVars(str string) string {
	re := regexp.MustCompile(`\$\{env:([^}]+)\}`)

	return re.ReplaceAllStringFunc(str, func(match string) string {
		envVar := strings.TrimPrefix(strings.TrimSuffix(match, "}"), "${env:")
		return os.Getenv(envVar)
	})
}
