package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return filepath.Join(home, configFileName), nil
}

func write(cfg Config) error {
	filePath, errGetPath := getConfigPath()
	if errGetPath != nil {
		return fmt.Errorf("failed to get config path: %w", errGetPath)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}
