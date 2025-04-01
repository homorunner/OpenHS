package bootstrap

import (
	"github.com/openhs/internal/card"
	"github.com/openhs/internal/config"
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// Initialize initializes all global components
func Initialize(configPath string) error {
	// Load global configuration first
	if err := config.LoadConfig(configPath); err != nil {
		return err
	}

	// Initialize logger with configuration
	logCfg := config.GetConfig().Log
	_, err := logger.SetLogger(logCfg.Level, logCfg.LogDir, logCfg.LogFile, logCfg.DevMode)
	if err != nil {
		return err
	}

	// Initialize card manager
	if err := card.InitializeCardManager(config.GetConfig().CardConfigDir); err != nil {
		return err
	}
	
	// Initialize game manager
	if err := game.InitializeGameManager(config.GetConfig().GameConfigDir); err != nil {
		return err
	}

	logger.Info("Global initialization complete")
	return nil
}
