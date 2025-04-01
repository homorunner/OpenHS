package config

import (
	"encoding/json"
	"os"
)

const (
	DefaultMaxMana          = 10
	DefaultStartingMana     = 0
	DefaultFirstPlayerHand  = 3
	DefaultSecondPlayerHand = 4
)

// GameConfig represents the configuration for a game
type GameConfig struct {
	Players []PlayerConfig `json:"players"`
}

// PlayerConfig represents the configuration for a player
type PlayerConfig struct {
	Hero string   `json:"hero"`
	Deck []string `json:"deck"`
}

// LoadGameConfig loads a game configuration from a JSON file
func LoadGameConfig(configPath string) (*GameConfig, error) {
	// Read the JSON file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data
	var config GameConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
