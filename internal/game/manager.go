package game

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/openhs/internal/logger"
)

var globalCardManager *CardManager

// GetCardManager returns the global card manager instance
func GetCardManager() *CardManager {
	if globalCardManager == nil {
		logger.Info("Initializing card manager")
		globalCardManager = NewCardManager()
	}
	return globalCardManager
}

// CardManager handles card creation and management
type CardManager struct {
	cardTemplates map[string]Card
}

// NewCardManager creates a new card manager
func NewCardManager() *CardManager {
	return &CardManager{
		cardTemplates: make(map[string]Card),
	}
}

// RegisterCard registers a new card template
func (cm *CardManager) RegisterCard(card Card) {
	logger.Debug("Registering card template", logger.String("name", card.Name))
	cm.cardTemplates[card.Name] = card
}

// CreateCardInstance creates a new instance of a card from a template
func (cm *CardManager) CreateCardInstance(name string) (*Card, error) {
	logger.Debug("Creating card instance", logger.String("name", name))
	template, exists := cm.cardTemplates[name]
	if !exists {
		err := NewCardError(ErrCardNotFound, fmt.Sprintf("card template not found: %s", name))
		logger.Error("Failed to create card", logger.String("name", name), logger.Err(err))
		return nil, err
	}

	card := template
	return &card, nil
}

// GetCardTemplate returns a card template by name
func (cm *CardManager) GetCardTemplate(name string) (*Card, error) {
	logger.Debug("Retrieving card template", logger.String("name", name))
	template, exists := cm.cardTemplates[name]
	if !exists {
		err := NewCardError(ErrCardNotFound, fmt.Sprintf("card template not found: %s", name))
		logger.Error("Failed to get card template", logger.String("name", name), logger.Err(err))
		return nil, err
	}
	return &template, nil
}

// GameManager handles loading and managing game configurations
type GameManager struct {
	games map[string]*GameConfig
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
		games: make(map[string]*GameConfig),
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
		gameConfig, err := LoadGameConfig(path)
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
func (gm *GameManager) GetGameConfig(gameID string) (*GameConfig, bool) {
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
