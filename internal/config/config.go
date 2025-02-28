package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Unable to get home directory: %s", err)
		return "", err
	}

	return fmt.Sprintf("%s/%s", homeDir, configFileName), nil
}

const configFileName = ".gatorconfig.json"

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

func Read() (*Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		log.Printf("Unable to get config file path: %s", err)
		return nil, err
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error opening gatorconfig in home directory: %s", err)
		return nil, err
	}

	config := Config{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %s", err)
		return nil, err
	}

	return &config, nil
}

func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		log.Printf("Unable to get config file path: %s", err)
		return err
	}

	configJson, err := json.Marshal(cfg)
	if err != nil {
		log.Printf("Unable to marshal json: %s", err)
		return err
	}

	return os.WriteFile(filePath, configJson, 0644)
}
