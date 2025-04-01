package game

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/openhs/internal/config"
	"github.com/openhs/internal/logger"
)

// GameManager handles loading and managing game configurations
type GameManager struct {
	games map[string]*config.GameConfig
}

var gameManager *GameManager

// GetGameManager returns the global game manager instance
func GetGameManager() *GameManager {
	if gameManager == nil {
		panic("game manager is not initialized")
	}
	return gameManager
}

// InitializeGameManager loads all game configurations from the specified directory
func InitializeGameManager(gameConfigDir string) error {
	gameManager = &GameManager{
		games: make(map[string]*config.GameConfig),
	}

	// Read all files in the game config directory
	err := filepath.WalkDir(gameConfigDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only process JSON files
		if filepath.Ext(path) != ".json" {
			return nil
		}

		// Load the game config
		gameConfig, err := config.LoadGameConfig(path)
		if err != nil {
			logger.Warn("Failed to load game config " + path + ": " + err.Error())
			return nil
		}

		// Use the filename without extension as the game ID
		gameID := filepath.Base(path[:len(path)-len(filepath.Ext(path))])
		gameManager.games[gameID] = gameConfig
		logger.Debug("Loaded game configuration: " + gameID)

		return nil
	})

	if err != nil {
		return err
	}

	logger.Info("Loaded " + fmt.Sprintf("%d", len(gameManager.games)) + " game configurations")
	return nil
}

// GetGameConfig returns a game configuration by ID
func (gm *GameManager) GetGameConfig(gameID string) (*config.GameConfig, bool) {
	config, exists := gm.games[gameID]
	return config, exists
}

// LoadGameByID loads a game by its configuration ID
func (gm *GameManager) LoadGameByID(gameID string) (*Game, error) {
	gameConfig, exists := gm.GetGameConfig(gameID)
	if !exists {
		return nil, fmt.Errorf("game configuration not found: %s", gameID)
	}
	
	return LoadGame(gameConfig)
} 