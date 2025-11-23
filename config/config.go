package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	LogLevel string `json:"log_level"`
}

func LoadConfig() (*Config, error) {
	// Default config
	config := &Config{
		LogLevel: "INFO",
	}

	file, err := os.Open("config.json")
	if os.IsNotExist(err) {
		// Create default config file if it doesn't exist
		return config, SaveConfig(config)
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func SaveConfig(config *Config) error {
	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
