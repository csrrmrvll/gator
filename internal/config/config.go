package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DB_url          string "json:\"db_url\""
	CurrentUserName string "json:\"current_user_name\""
}

func Read() Config {
	configFile := getConfigFilePath()
	if configFile == "" {
		return Config{}
	}
	jsonFile, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}
	}

	var config Config
	err = json.Unmarshal(jsonFile, &config)
	if err != nil {
		return Config{}
	}
	return config
}

func SetUser(config Config) {
	config.CurrentUserName = os.Getenv("USER")
	write(config)
}

func getConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homeDir + "/" + configFileName
}
func write(config Config) {
	configFile := getConfigFilePath()
	jsonData, err := json.Marshal(config)
	if err != nil {
		return
	}
	err = os.WriteFile(configFile, jsonData, 0644)
	if err != nil {
		return
	}
}
