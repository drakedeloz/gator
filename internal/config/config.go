package config

import (
	"encoding/json"
	"log"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_URL      string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func Read() *Config {
	var cfg Config
	configPath, err := getConfigFilePath()
	if err != nil {
		return &Config{}
	}
	fileBytes, _ := os.ReadFile(configPath)
	err = json.Unmarshal(fileBytes, &cfg)
	if err != nil {
		log.Fatalf("could not read config file: %v", err)
		return &Config{}
	}

	return &cfg
}

func (m *Config) SetUser(user string) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return
	}
	m.CurrentUser = user
	configBytes, err := json.Marshal(m)
	if err != nil {
		log.Fatalf("could not marshal config: %v", err)
		return
	}

	err = os.WriteFile(configPath, configBytes, 0644)
	if err != nil {
		log.Fatalf("could not write to config file: %v", err)
		return
	}
}

func getConfigFilePath() (string, error) {
	configPath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get user")
		return "", err
	}
	return configPath + "/" + configFileName, nil
}
