package config

import (
	"encoding/json"
	"os"
)

const JsonFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	homeDir, _ := os.UserHomeDir()
	path := homeDir + "/" + JsonFileName
	jsonFile, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	defer jsonFile.Close()

	var cfg Config
	decoder := json.NewDecoder(jsonFile)
	if err = decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func wrtie(cfg *Config) error {
	homeDir, _ := os.UserHomeDir()
	path := homeDir + "/" + JsonFileName

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	return wrtie(cfg)
}
