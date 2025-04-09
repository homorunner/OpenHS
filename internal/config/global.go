package config

import (
	"encoding/json"
	"os"
)

// LogConfig represents the logging configuration
type LogConfig struct {
	Level   string `json:"level"`
	LogDir  string `json:"log_dir"`
	LogFile string `json:"log_file"`
	DevMode bool   `json:"dev_mode"`
}

// GlobalConfig represents the global configuration
type GlobalConfig struct {
	GameConfigDir string    `json:"game_config_dir"`
	Log           LogConfig `json:"logging"`
}

var globalConfig *GlobalConfig

// GetConfig returns the global configuration instance
func GetConfig() *GlobalConfig {
	if globalConfig == nil {
		panic("global config is not initialized")
	}
	return globalConfig
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config GlobalConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	globalConfig = &config
	return nil
}
