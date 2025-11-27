package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	config := Config{}

	filePath, errGetPath := getConfigPath()
	if errGetPath != nil {
		return config, fmt.Errorf("failed to get config path: %w", errGetPath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("failed to decode config: %w", err)
	}

	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username

	errWrite := write(*c)
	if errWrite != nil {
		return fmt.Errorf("failed to write config file: %w", errWrite)
	}

	return nil
}
