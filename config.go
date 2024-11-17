package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const configPath = "/.config/tempus/config.json"

type config struct {
	WebhookUrl string `json:"webhookUrl"`
}

func readConfig() (config, error) {
	var cfg config
	home, _ := os.UserHomeDir()
	f, err := os.Open(home + configPath)
	if err != nil {
		return cfg, fmt.Errorf("error reading config: %w", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("error decoding config: %w", err)
	}

	return cfg, nil
}

func writeConfig(cfg config) error {

	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}

	err = json.NewEncoder(f).Encode(&cfg)
	if err != nil {
		return fmt.Errorf("error writing config: %w", err)
	}
	return nil
}
