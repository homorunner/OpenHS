package game

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/openhs/internal/logger"
)

var globalCardManager *CardManager

// GetCardManager returns the global card manager instance
func GetCardManager() *CardManager {
	if globalCardManager == nil {
		globalCardManager = NewCardManager()
	}
	return globalCardManager
}

// InitializeCardManager initializes the global card manager and loads card database
func InitializeCardManager(configDir string) error {
	cm := GetCardManager()
	return cm.LoadCardDatabase(configDir)
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

// CreateCard creates a new instance of a card from a template
func (cm *CardManager) CreateCard(name string) (*Card, error) {
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

// LoadCardDatabase loads all card templates from the specified directory
func (cm *CardManager) LoadCardDatabase(cardConfigDir string) error {
	logger.Info("Loading card database", logger.String("dir", cardConfigDir))

	// Ensure the card config directory exists
	if _, err := os.Stat(cardConfigDir); os.IsNotExist(err) {
		err := NewCardError(ErrCardNotFound, fmt.Sprintf("card config directory not found: %s", cardConfigDir))
		logger.Error("Failed to load card database", logger.String("dir", cardConfigDir), logger.Err(err))
		return err
	}

	// Read all JSON files in the cards directory
	files, err := os.ReadDir(cardConfigDir)
	if err != nil {
		logger.Error("Failed to read card config directory", logger.String("dir", cardConfigDir), logger.Err(err))
		return err
	}

	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(cardConfigDir, file.Name()))
		if err != nil {
			logger.Error("Failed to read card config file",
				logger.String("file", file.Name()),
				logger.Err(err))
			continue
		}

		var cardConfig CardConfig
		if err := json.Unmarshal(data, &cardConfig); err != nil {
			logger.Error("Failed to parse card config",
				logger.String("file", file.Name()),
				logger.Err(err))
			continue
		}

		// Convert config to Card
		card := Card{
			Name:      cardConfig.Name,
			Cost:      cardConfig.Cost,
			Attack:    cardConfig.Attack,
			Health:    cardConfig.Health,
			MaxHealth: cardConfig.Health, // Set MaxHealth equal to Health
			Type:      cardConfig.Type,
		}

		// Register the card template
		cm.RegisterCard(card)
		logger.Info("Loaded card template",
			logger.String("name", card.Name),
			logger.String("file", file.Name()))
	}

	logger.Info("Card database loaded successfully")
	return nil
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
