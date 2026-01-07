package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type config struct {
	StartLines int `json:"startLines"`
}

var Cfg *config = &config{}

func LoadConfig() error {
	err := ensureConfigFile()
	if err != nil {
		return err
	}

	configDir, err := getConfigDir()

	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "runner.json")

	conf := config{
		StartLines: 20,
	}

	data, err := os.ReadFile(configPath)
	if err = json.Unmarshal(data, &conf); err != nil {
		return err
	}
	Cfg = &conf
	return nil
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".config", "runner")
	return configDir, nil
}

func GetConfig() (*config, error) {
	return nil, nil
}

func ensureAndGetConfigDirectory() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return configDir, nil
}

func createConfigFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	b, err := json.Marshal(config{StartLines: 20})
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func ensureConfigFile() error {
	configDir, err := ensureAndGetConfigDirectory()
	if err != nil {
		return err
	}
	configPath := filepath.Join(configDir, "runner.json")
	if _, err = os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		if err = createConfigFile(configPath); err != nil {
			return err
		}
	}

	return nil
}
